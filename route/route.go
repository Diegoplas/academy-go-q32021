package route

import (
	"net/http"

	"github.com/Diegoplas/go-bootcamp-deliverable/controller"
	"github.com/Diegoplas/go-bootcamp-deliverable/csvdata"
	"github.com/Diegoplas/go-bootcamp-deliverable/service"

	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {

	pokemonService := service.NewGetPokemonService(csvdata.PokemonRepo{})
	pokemonHandler := controller.NewGetPokemonHandler(pokemonService)

	router = mux.NewRouter()
	router.HandleFunc("/first-generation/{id}", pokemonHandler.GetPokemonFromCSVHandler).Methods(http.MethodGet)
	router.HandleFunc("/second-generation/{id}", controller.GetPokemonExternalAPIHandler).Methods(http.MethodGet)
	return
}
