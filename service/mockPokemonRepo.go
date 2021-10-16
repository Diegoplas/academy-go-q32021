package service

import "github.com/Diegoplas/go-bootcamp-deliverable/model"

var MockedPokemonResponse = []model.PokemonData{
	{
		ID:     94,
		Name:   "Gengar",
		Height: 15,
		Type1:  "Ghost",
		Type2:  "Poison",
	},
	{
		ID:     95,
		Name:   "Onix",
		Height: 88,
		Type1:  "Rock",
		Type2:  "Ground",
	},
	{
		ID:     96,
		Name:   "Drowzee",
		Height: 10,
		Type1:  "Psychic",
		Type2:  " - ",
	},
}
