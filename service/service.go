package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/Diegoplas/go-bootcamp-deliverable/model"
)

type getter interface {
	ListPokemons() ([]model.PokemonData, error)
}

type PokemonRepo struct {
	repository getter
}

func NewRepositoryService(repository getter) PokemonRepo {
	return PokemonRepo{repository: repository}
}

func (pr PokemonRepo) GetPokemonFromCSV(wantedIndex string) (model.PokemonData, error) {

	allPokemons, err := pr.repository.ListPokemons()
	if err != nil {
		fmt.Printf("Error listing pokemons %s\n", err)
	}

	requestedIndex, err := strconv.Atoi(wantedIndex)
	if err != nil {
		log.Println("Error converting string to int")
		return model.PokemonData{}, fmt.Errorf("error: something happened")
	}

	for _, pokemon := range allPokemons {
		if pokemon.ID == requestedIndex {
			return pokemon, nil
		}
	}

	return model.PokemonData{}, fmt.Errorf("error: no pokemon found")
}

func (pr PokemonRepo) GetPokemonFromExternalAPI(wantedIndex string) (model.PokemonData, error) {

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s?name", wantedIndex)

	response, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return model.PokemonData{}, fmt.Errorf("request error")
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error reading response body: %v\n", err.Error())
		return model.PokemonData{}, fmt.Errorf("error reading data")
	}

	pokemonData := model.PokemonExternalData{}
	unmarshalErr := json.Unmarshal(data, &pokemonData)

	if unmarshalErr != nil {
		log.Printf("unmarshal failed: %v\n", err.Error())
		return model.PokemonData{}, fmt.Errorf("error getting data %v", err.Error())
	}

	formatedPokemonData := formatPokemonData(pokemonData)

	return formatedPokemonData, nil
}

func formatPokemonData(externalPokemonData model.PokemonExternalData) model.PokemonData {

	formatedPokemonData := model.PokemonData{}

	if len(externalPokemonData.Types) == 1 {
		formatedPokemonData = model.PokemonData{
			ID:     externalPokemonData.ID,
			Name:   externalPokemonData.Name,
			Height: externalPokemonData.Height,
			Type1:  externalPokemonData.Types[0].Type.Name,
			Type2:  " - ",
		}
	} else if len(externalPokemonData.Types) == 2 {
		formatedPokemonData = model.PokemonData{
			ID:     externalPokemonData.ID,
			Name:   externalPokemonData.Name,
			Height: externalPokemonData.Height,
			Type1:  externalPokemonData.Types[0].Type.Name,
			Type2:  externalPokemonData.Types[1].Type.Name,
		}
	}

	return formatedPokemonData
}
