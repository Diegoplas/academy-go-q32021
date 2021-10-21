package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/Diegoplas/go-bootcamp-deliverable/model"
)

type mockGetter struct {
	listPokemonResponse []model.PokemonData
	wantErr             bool
}

func (mg mockGetter) ListPokemons() ([]model.PokemonData, error) {
	if mg.wantErr {
		return nil, fmt.Errorf("list pokemons error")
	}
	return mg.listPokemonResponse, nil
}

func TestGetPokemonService_GetPokemonFromCSV(t *testing.T) {

	tests := []struct {
		name         string
		mockedGetter mockGetter
		requestedID  string
	}{
		{
			name:         "Valid test",
			mockedGetter: mockGetter{listPokemonResponse: MockedPokemonResponse},
			requestedID:  "95",
		},
		{
			name:         "Invalid - empty param",
			mockedGetter: mockGetter{wantErr: true},
			requestedID:  "",
		},
		{
			name:         "Invalid - error on CSVdata functions",
			mockedGetter: mockGetter{wantErr: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRepo := NewRepositoryService(tt.mockedGetter)
			gotPokemon, err := testRepo.GetPokemonFromCSV(tt.requestedID)
			if (err != nil) != tt.mockedGetter.wantErr {
				t.Errorf("GetPokemonFromCSV() error = %v, wantErr %v", err, tt.mockedGetter.wantErr)
				return
			}
			gotID := ""
			if gotPokemon.ID != 0 {
				gotID = strconv.Itoa(gotPokemon.ID)
			}
			if gotID != tt.requestedID {
				t.Errorf("GetPokemonFromCSV() Got ID = %v, wanted ID %v", gotID, tt.requestedID)
			}
		})
	}
}

type mockRequestSender struct {
	method   string
	url      string
	response *http.Response
	wantErr  bool
}

type mockRequestParser struct {
	wantErr bool
}

func (mrs *mockRequestSender) SendRequest(method, url, values map[string]interface{}) (*http.Response, error) {
	if mrs.wantErr {
		return nil, fmt.Errorf("send request error")
	}
	return mrs.response, nil
}

func (mrp *mockRequestParser) ParseResponse(resp *http.Response) ([]byte, error) {
	if mrp.wantErr {
		return nil, fmt.Errorf("request parser error")
	}
	parsedResponse, _ := ioutil.ReadAll(resp.Body)
	return parsedResponse, nil
}

func init() {
	Client = &mockClient{}
}

type mockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func mockClientDoFactory(func(req *http.Request) (*http.Response, error)) mockClient {
	return mockClient{DoFunc: GetDoFunc}
}
func TestGetPokemonService_GetPokemonFromExternalAPI(t *testing.T) {

	externalJSONResponse := `{
		"id": 251,
		"name": "celebi",
		"height": 6,
		"types": [
			{
				"slot": 1,
				"type": {
					"name": "psychic",
					"url": "https://pokeapi.co/api/v2/type/14/"
				}
			},
			{
				"slot": 2,
				"type": {
					"name": "grass",
					"url": "https://pokeapi.co/api/v2/type/12/"
				}
			}
		],
	}`

	formatedResponse := model.PokemonData{
		ID:     251,
		Name:   "celebi",
		Height: 6,
		Type1:  "psychic",
		Type2:  "grass",
	}

	validRequestedID := "251"
	tests := []struct {
		name           string
		requestedID    string
		mockClient     mockClient
		wantedResponse model.PokemonData
		wantErr        bool
	}{
		{
			name:        "Valid test",
			requestedID: validRequestedID,
			mockClient: mockClient{func(*http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(bytes.NewBufferString(externalJSONResponse)),
				}, nil
			}},
			wantedResponse: formatedResponse,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testMocker := mockClientDoFactory(GetDoFunc)
			gotPokemon, err := testMocker.GetPokemonFromExternalAPI(tt.requestedID)
			if gotPokemon != tt.wantedResponse {
				t.Errorf("GetPokemonFromExternalAPI() Got ID = %v, wanted ID %v", gotPokemon, tt.wantedResponse)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPokemonFromExternalAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
