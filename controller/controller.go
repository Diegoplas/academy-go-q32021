package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Diegoplas/go-bootcamp-deliverable/csvdata"
	"github.com/Diegoplas/go-bootcamp-deliverable/model"
	"github.com/Diegoplas/go-bootcamp-deliverable/service"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

func GetPokemonFromCSVHandler(w http.ResponseWriter, r *http.Request) {

	requestIndex := mux.Vars(r)["id"]

	wantedIndex, err := validateFirstGenID(requestIndex)
	if err != nil {
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	response, err := service.GetPokemonFromCSV(wantedIndex)
	if err != nil {
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusInternalServerError, errorResponse)
		return
	}

	render.New().JSON(w, http.StatusOK, &response)
}

func GetPokemonExternalAPIHandler(w http.ResponseWriter, r *http.Request) {

	requestIndex := mux.Vars(r)["id"]
	validatedIndex, err := validateSecondGenID(requestIndex)
	if err != nil {
		errorResponse := model.ErrorResponse{Err: err.Error()}
		render.New().JSON(w, http.StatusBadRequest, errorResponse)
		return
	}

	var response model.PokemonData

	response, err = service.GetPokemonFromExternalAPI(validatedIndex)
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
