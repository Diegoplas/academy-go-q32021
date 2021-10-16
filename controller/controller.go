package controller

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Diegoplas/go-bootcamp-deliverable/config"
	"github.com/Diegoplas/go-bootcamp-deliverable/csvdata"
	"github.com/Diegoplas/go-bootcamp-deliverable/model"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

type getPokemonInfo interface {
	GetPokemonFromCSV(wantedIndex string) (model.PokemonData, error)
	GetPokemonFromExternalAPI(wantedIndex string) (model.PokemonData, error)
}

type PokemonHandler struct {
	getPokemons getPokemonInfo
}

func NewGetPokemonHandler(getPokemons getPokemonInfo) PokemonHandler {
	return PokemonHandler{getPokemons: getPokemons}
}

func (pk PokemonHandler) GetPokemonFromCSVHandler(w http.ResponseWriter, r *http.Request) {

	requestIndex := mux.Vars(r)["id"]

	wantedIndex, err := validateFirstGenID(requestIndex)
	if err != nil {
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	response, err := pk.getPokemons.GetPokemonFromCSV(wantedIndex)
	if err != nil {
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusInternalServerError, errorResponse)
		return
	}

	render.New().JSON(w, http.StatusOK, &response)
}

func validateFirstGenID(index string) (string, error) {

	wantedIndex, err := strconv.Atoi(index)
	if err != nil {
		log.Printf("string to int convertion failed %v", err.Error())
		return "", fmt.Errorf("something happened")
	}

	if wantedIndex < 1 || wantedIndex > 152 {
		log.Println("invalid id")
		return "", errors.New("please introduce a valid pokemon ID from First gen. (1-151)")

	}

	return index, nil
}

func (pk PokemonHandler) GetPokemonExternalAPIHandler(w http.ResponseWriter, r *http.Request) {

	requestIndex := mux.Vars(r)["id"]
	validatedIndex, err := validateSecondGenID(requestIndex)
	if err != nil {
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	response, err := pk.getPokemons.GetPokemonFromExternalAPI(validatedIndex)
	if err != nil {
		log.Printf("Error on GetPokemonFromExternalAPI %s", err.Error())
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusInternalServerError, errorResponse)
		return
	}
	err = csvdata.WritePokemonIntoCSV(response)
	if err != nil {
		log.Printf("Error on WritePokemonIntoCSV: %s", err.Error())
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusInternalServerError, errorResponse)
		return
	}

	render.New().JSON(w, http.StatusOK, &response)
}

func validateSecondGenID(index string) (string, error) {

	wantedIndex, err := strconv.Atoi(index)
	if err != nil {
		log.Printf("string to int convertion failed %v", err.Error())
		return "", fmt.Errorf("something happened")
	}

	if wantedIndex < 152 || wantedIndex > 251 {
		log.Println("invalid id")
		return "", errors.New("please introduce a valid pokemon ID from Second gen. (152-251)")
	}

	return index, nil
}

func worker(readChan chan []string, writeChan chan []string) {
	for {
		select {
		case row, ok := <-writeChan:
			if !ok {
				return
			}
			readChan <- row
		}
	}
}

func WorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
	numberType := strings.ToLower(r.URL.Query().Get("type"))
	if (numberType != "odd" && numberType != "even") || numberType == "" {
		errorResponse := model.ErrorResponse{Err: "Invalid or empty param Type."}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	items, err := strconv.Atoi(r.URL.Query().Get("items"))
	if err != nil || items == 0 {
		log.Println("Error converting items param to int.")
		errorResponse := model.ErrorResponse{Err: "Invalid or empty items param."}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	itemsPerWorker, err := strconv.Atoi(r.URL.Query().Get("items_per_workers"))
	if err != nil || itemsPerWorker == 0 {
		log.Println("Error converting items_per_workers param to int.")
		errorResponse := model.ErrorResponse{Err: "Invalid or empty items per worker param."}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	if itemsPerWorker > items {
		errorResponse := model.ErrorResponse{Err: "Items per worker param should be lower than items param."}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
	}

	numberOfWorkers := items / itemsPerWorker

	csvfile, err := os.Open(config.FirstGenCSVPath)
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	// create the input/output channels
	srcCh := make(chan []string)
	outputCh := make(chan []string, items)

	// manage synchronization
	var wg sync.WaitGroup

	// declare the workers
	for workerID := 1; workerID <= numberOfWorkers; workerID++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(outputCh, srcCh)
		}()
	}

	rowNumber := 0
	itemsToGet := items * 2 //total items to get, after odd or even choice
	go func() {
		for {
			record, err := reader.Read()

			if err == io.EOF {
				fmt.Println("ending because of EOF")
				break
			} else if err != nil {
				log.Println(err)
				break
			}
			rowNumber += 1
			if rowNumber <= itemsToGet {
				idNum, err := strconv.Atoi(record[0])
				if err != nil {
					log.Println(err)
					break
				}
				if numberType == "even" && (idNum%2 == 0) {
					srcCh <- record
				} else if numberType == "odd" && (idNum%2 == 1) {
					srcCh <- record
				}
			} else {
				break
			}
		}
		close(srcCh)
	}()

	// wait for worker group to finish.
	go func() {
		wg.Wait()
		close(outputCh)
	}()

	// format the response.
	response := []model.PokemonData{}
	for pokemonInfo := range outputCh {
		pokemonID, err := strconv.Atoi(pokemonInfo[0])
		if err != nil {
			log.Println(err)
			break
		}
		pokemonHeight, err := strconv.Atoi(pokemonInfo[2])
		if err != nil {
			log.Println(err)
			break
		}
		pokemon := model.PokemonData{
			ID:     pokemonID,
			Name:   pokemonInfo[1],
			Height: pokemonHeight,
			Type1:  pokemonInfo[3],
			Type2:  pokemonInfo[4],
		}
		response = append(response, pokemon)

	}

	render.New().JSON(w, http.StatusOK, &response)
}
