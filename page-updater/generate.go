package main

import (
	"bufio"
	"fmt"
	"strings"
	"text/template"

	"github.com/taichi765/kokkimusume-wiki-automation/types"
)

func generateCharaListPage(old string, charas []types.CharacterData) (string, error) {
	beforeStart, afterEnd, err := splitLinesToEdit(old)
	if err != nil {
		return "", err
	}

	generated, err := generateCharaListPageContent(charas)
	if err != nil {
		return "", fmt.Errorf("failed to generate content for CharaList page: %w", err)
	}

	return beforeStart + generated + afterEnd, nil
}

// Part of [generateCharaListPage].
func generateCharaListPageContent(charas []types.CharacterData) (string, error) {
	tmpl, err := template.New("table_item").Parse(
		`|[[{{.Name}}]]|{{.Area}}|{{.FirstAppearenceDate}}|
`)
	if err != nil {
		panic("template must be valid")
	}

	b := &strings.Builder{}
	fmt.Fprintln(b, "#tablesort{{")
	fmt.Fprintln(b, "|名前|地域|初出|h")
	for _, ch := range charas {
		err := tmpl.Execute(b, ch)
		if err != nil {
			return "", fmt.Errorf("failed to generate table item: %w", err)
		}
	}
	fmt.Fprintln(b, "}}")

	return b.String(), nil
}

func generateMenuBar(old string, charas []types.CharacterData) (string, error) {
	beforeStart, afterEnd, err := splitLinesToEdit(old)
	if err != nil {
		return "", err
	}

	generated := generateMenuBarContent(charas)

	return beforeStart + generated + afterEnd, nil
}

// Part of [generateMenuBar].
func generateMenuBarContent(charas []types.CharacterData) string {
	byArea := charasByArea(charas)
	b := &strings.Builder{}
	for _, a := range types.ValidAreas {
		chs := byArea[a]
		fmt.Fprintf(b, "** %v\n", a)
		for _, c := range chs {
			fmt.Fprintf(b, "- [[%v]]\n", c.Name)
		}
	}
	return b.String()
}

func charasByArea(charas []types.CharacterData) map[string][]types.CharacterData {
	byArea := make(map[string][]types.CharacterData, len(types.ValidAreas))
	for _, ch := range charas {
		if charas, ok := byArea[ch.Area]; ok {
			byArea[ch.Area] = append(charas, ch)
		} else {
			byArea[ch.Area] = []types.CharacterData{ch}
		}
	}
	return byArea
}

// @generated_startがある行まで(含む)と@generated_endがある行以降(含む)
func splitLinesToEdit(old string) (string, string, error) {
	sc := bufio.NewScanner(strings.NewReader(old))

	var beforeStart, afterEnd strings.Builder
	// 0: finding @generated_start
	// 1: finding @generated_end
	// 2: after @generated_end
	phase := 0
	for sc.Scan() {
		line := sc.Text()
		trimmed := strings.Join(strings.Fields(line), "")

		if strings.HasPrefix(trimmed, "//@generated_end") {
			switch phase {
			case 0:
				return "", "", fmt.Errorf("found @generated_end before finding @generated_start")
			case 1:
				phase = 2
			case 2:
				return "", "", fmt.Errorf("multiple @generated_end was found")
			default:
				panic("invalid phase")
			}
		}

		switch phase {
		case 0:
			_, err := beforeStart.WriteString(line + "\n")
			if err != nil {
				return "", "", err
			}
		case 1:
			// do nothing
		case 2:
			_, err := afterEnd.WriteString(line + "\n")
			if err != nil {
				return "", "", err
			}
		default:
			panic("invalid phase")
		}

		if strings.HasPrefix(trimmed, "//@generated_start") {
			if phase != 0 {
				return "", "", fmt.Errorf("multiple `@generated_start` was found")
			}
			phase = 1
		}
	}

	if err := sc.Err(); err != nil {
		return "", "", fmt.Errorf("an error occured while reading old page content: %w", err)
	}

	if phase != 2 {
		return "", "", fmt.Errorf("can't find @generated_start or @generated_end")
	}

	return beforeStart.String(), afterEnd.String(), nil
}
