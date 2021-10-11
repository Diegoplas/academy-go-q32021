package csvdata

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Diegoplas/go-bootcamp-deliverable/config"
	"github.com/Diegoplas/go-bootcamp-deliverable/model"
)

type PokemonRepo struct {
}

func (pr PokemonRepo) ListPokemons() ([]model.PokemonData, error) {

	// open the file
	csvFile, err := os.Open(config.FirstGenCSVPath)
	if err != nil {
		log.Printf("Error opening csv file %v:", err.Error())
		return nil, fmt.Errorf("error loading database")
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		log.Printf("Error csv reader %v:", err.Error())
		return nil, fmt.Errorf("error reading database")
	}

	allPokemons, err := linesToSlice(csvLines)
	if err != nil {
		log.Printf("Error opening csv file %v:", err.Error())
		return nil, fmt.Errorf("data handling error")
	}

	return allPokemons, nil
}

func linesToSlice(csvLines [][]string) ([]model.PokemonData, error) {

	totalPokemonsFirstRegion := 151
	allPokemons := make([]model.PokemonData, totalPokemonsFirstRegion)

	for _, line := range csvLines {

		pokemonID, err := strconv.Atoi(line[0])
		if err != nil {
			return nil, fmt.Errorf("error converting string to int %v", err.Error())
		}

		pokemonHeight, err := strconv.Atoi(line[2])
		if err != nil {
			return nil, fmt.Errorf("error converting string to int %v", err.Error())
		}

		if line[4] == "" {
			line[4] = " - "
		}

		pokemon := model.PokemonData{
			ID:     pokemonID,
			Name:   line[1],
			Height: pokemonHeight,
			Type1:  line[3],
			Type2:  line[4],
		}
		allPokemons = append(allPokemons, pokemon)
	}

	return allPokemons, nil
}

func WritePokemonIntoCSV(externalPokemonData model.PokemonData) error {

	strPokemonID := strconv.Itoa(externalPokemonData.ID)
	row := []string{

		strPokemonID,
		strings.Title(externalPokemonData.Name),
		strconv.Itoa(externalPokemonData.Height),
		externalPokemonData.Type1,
		externalPokemonData.Type2,
	}

	csvFile, err := os.OpenFile(config.SecondGenCSVPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Printf("Error opening csv file: %s", err)
		return fmt.Errorf("database error")
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return fmt.Errorf("error reading csv file: %v", err.Error())
	}

	for _, line := range csvLines {
		if line[0] == strPokemonID {
			return nil
		}
	}

	writer := csv.NewWriter(csvFile)

	errWrite := writer.Write(row)
	if errWrite != nil {
		log.Println("Error writing into csv file:", errWrite)
	}
	defer writer.Flush()

	return nil
}
