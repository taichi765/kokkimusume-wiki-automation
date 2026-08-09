package main

import (
	"bufio"
	"fmt"
	"strings"
	"text/template"
)

func editCharaListPage(old string, charas []CharacterData) (string, error) {
	start, end, err := findLinesToEdit(old)
	if err != nil {
		return "", err
	}

	generated, err := generateCharaListPageContent(charas)
	if err != nil {
		return "", fmt.Errorf("failed to generate content for CharaList page: %w", err)
	}

	return old[:start+1] + generated + old[end:], nil
}

func generateCharaListPageContent(charas []CharacterData) (string, error) {
	tmpl, err := template.New("item").Parse(`
	|{{.name}}|{{.area}}|{{.firstAppearenceDate}}|
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

func editMenuBar(old string) (string, error) {
	panic("TODO")
}

// Finds lines to be edited, which starts by `@generated_start` and ends by `@generated_end`.
// Returns `(-1, -1, <error>)` if an error occured.
//
// Line number starts with 0.
func findLinesToEdit(old string) (int, int, error) {
	sc := bufio.NewScanner(strings.NewReader(old))

	cnt := 0
	var start_line int
	start_line_found := false
	var end_line int
	end_line_found := false
	for sc.Scan() {
		line := sc.Text()
		line = strings.Join(strings.Fields(line), "")

		if strings.HasPrefix(line, "//@generated_start") {
			if start_line_found {
				return -1, -1, fmt.Errorf("multiple `@generated_start` was found: first at %v, secound at %v", start_line, cnt)
			}

			start_line = cnt
			start_line_found = true
		}

		if strings.HasPrefix(line, "//@generated_end") {
			if end_line_found {
				return -1, -1, fmt.Errorf("multiple `@generated_end` was found: first at %v, secound at %v", end_line, cnt)
			}
			if !start_line_found {
				return -1, -1, fmt.Errorf("found `@generated_end` at %v before finding `@generated_start`", cnt)
			}

			end_line = cnt
			end_line_found = true
		}

		cnt += 1
	}

	if err := sc.Err(); err != nil {
		return -1, -1, fmt.Errorf("an error occured while reading old page content: %w", err)
	}

	if !start_line_found {
		return -1, -1, fmt.Errorf("cannot find line starts with @generated_start")
	}
	if !end_line_found {
		return -1, -1, fmt.Errorf("cannot find line starts with `@generated_end`")
	}

	return start_line, end_line, nil
}
