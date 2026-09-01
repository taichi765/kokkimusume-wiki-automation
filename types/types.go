package types

import "slices"

type CharacterData struct {
	Name                string
	Area                string
	FirstAppearenceDate string
}

var ValidAreas = []string{"東アジア", "東南アジア・南アジア", "中東", "ヨーロッパ", "オセアニア", "北米", "中南米", "アフリカ"}

// Returns whether the given area is valid or not.
func AreaIsValid(area string) bool {
	return slices.Contains(ValidAreas, area)
}
