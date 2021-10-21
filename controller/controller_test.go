package controller

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Diegoplas/go-bootcamp-deliverable/model"
	"github.com/gorilla/mux"
)

type handlersMock struct {
	CSVResponse      model.PokemonData
	externalResponse model.PokemonData
	csvReader        *csv.Reader
	wantErr          bool
}

func (hm handlersMock) GetPokemonFromCSV(wantedIndex string) (model.PokemonData, error) {
	if hm.wantErr {
		return model.PokemonData{}, fmt.Errorf("error getting pokemon from csv file")
	}
	return hm.CSVResponse, nil
}

func (hm handlersMock) GetPokemonFromExternalAPI(wantedIndex string) (model.PokemonData, error) {
	if hm.wantErr {
		return model.PokemonData{}, fmt.Errorf("error getting pokemon from external api")
	}
	return hm.externalResponse, nil
}

func (hm handlersMock) CreateReaderFromCSVFile(csvPath string) (*csv.Reader, error) {
	if hm.wantErr {
		return nil, fmt.Errorf("error making csv reader")
	}
	return hm.csvReader, nil
}

func TestPokemonHandler_GetPokemonFromCSVHandler(t *testing.T) {

	endpoint := "/first-generation/{id}"

	validResponse := model.PokemonData{
		ID:     94,
		Name:   "Gengar",
		Height: 15,
		Type1:  "Ghost",
		Type2:  "Poison",
	}

	tests := []struct {
		name          string
		wantCode      int
		mockedService handlersMock
		requestedID   string
	}{
		{
			name:        "Valid test",
			requestedID: "94",
			mockedService: handlersMock{
				CSVResponse: validResponse,
			},
			wantCode: http.StatusOK,
		},
		{
			name:          "Invalid - empty param",
			requestedID:   "",
			mockedService: handlersMock{},
			wantCode:      http.StatusBadRequest,
		},
		{
			name:        "Invalid - get pokemon from csv error",
			requestedID: "94",
			mockedService: handlersMock{
				wantErr: true,
			},
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, endpoint, nil)
			r = mux.SetURLVars(r, map[string]string{"id": tt.requestedID})
			testPokemonHandler := NewGetPokemonHandler(tt.mockedService)
			testPokemonHandler.GetPokemonFromCSVHandler(w, r)
			if w.Code != tt.wantCode {
				t.Errorf("Error on GetPokemonFromCSVHandler: got code: %v expected code: %v", w.Code, tt.wantCode)
				return
			}
		})
	}
}
