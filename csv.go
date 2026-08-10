package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
)

const CSV_PATH = "./data/charas.csv"

// Loads character data from csv.
func loadCharaData() ([]CharacterData, error) {
	f, err := os.Open(CSV_PATH)
	if err != nil {
		return nil, fmt.Errorf("failed to open csv file: %w", err)
	}
	defer f.Close()

	rd := csv.NewReader(f)

	// validate header
	r, err := rd.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read from csv file")
	}
	if r[0] != "Name" || r[1] != "Area" || r[2] != "Date" {
		return nil, fmt.Errorf("invalid header: %v", r)
	}

	charas := []CharacterData{}
	for {
		r, err := rd.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if !areaIsValid(r[1]) {
			return nil, fmt.Errorf("invalid area: %v", r[1])
		}
		charas = append(charas, CharacterData{
			Name:                r[0],
			Area:                r[1],
			FirstAppearenceDate: r[2],
		})
	}

	return charas, nil
}

// Returns whether the given area is valid or not.
func areaIsValid(area string) bool {
	validAreas := []string{"東アジア", "東南アジア・南アジア", "中東", "ヨーロッパ", "オセアニア", "北米", "中南米", "アフリカ"}
	return slices.Contains(validAreas, area)
}
