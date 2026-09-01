package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/taichi765/kokkimusume-wiki-automation/types"
)

const CSV_PATH = "../data/charas.csv"

// Loads character data from csv.
func loadCharaData() ([]types.CharacterData, error) {
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

	charas := []types.CharacterData{}
	for {
		r, err := rd.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if !types.AreaIsValid(r[1]) {
			return nil, fmt.Errorf("invalid area: %v", r[1])
		}
		charas = append(charas, types.CharacterData{
			Name:                r[0],
			Area:                r[1],
			FirstAppearenceDate: r[2],
		})
	}

	return charas, nil
}
