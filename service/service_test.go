package service

import (
	"reflect"
	"testing"

	"github.com/Diegoplas/go-bootcamp-deliverable/model"
)

// var wantedResponse = model.PokemonData{
// 	ID:   94,
// 	Name: "Gengar",
// }

func TestGetPokemonService_GetPokemonFromCSV(t *testing.T) {

	type mockGetter struct {
		repository getter
		listPokemonResponse []model.PokemonData
		wantErr             bool
	}

	type getter interface {
		ListPokemons() ([]model.PokemonData, error)
	}

	func (mg mockGetter) ListPokemons()([]model.PokemonData, error) {
		if mg.wantErr{
			return nil, fmt.Errorf("list pokemons error")
		}
		return mg.listPokemonResponse, nil
	}

	type fields struct {
		listPokemonRepo getter
	}

	type args struct {
		wantedIndex string
	}

	tests := []struct {
		name     string
		fields   fields
		response model.PokemonData
		wantErr  bool
		args     args
	}{
		{
			name: "Pokemon obtained correctly",
			fields: fields{
				listPokemonRepo: mockGetter{listPokemonResponse: MockedPokemonRepo},
			},
			response: wantedResponse,
			wantErr:  false,
		},
	}
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, err := tt.gps.repository.GetPokemonFromCSV(tt.args.wantedIndex)
	// 		if (err != nil) != tt.wantErr {
	// 			t.Errorf("GetPokemonService.GetPokemonFromCSV() error = %v, wantErr %v", err, tt.wantErr)
	// 			return
	// 		}
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("GetPokemonService.GetPokemonFromCSV() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
