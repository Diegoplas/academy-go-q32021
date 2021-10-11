package route

import (
	"net/http"

	"github.com/Diegoplas/go-bootcamp-deliverable/controller"

	"github.com/gorilla/mux"
)

func GetRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/first-generation/{id}", controller.GetPokemonFromCSVHandler).Methods(http.MethodGet)
	router.HandleFunc("/second-generation/{id}", controller.GetPokemonExternalAPIHandler).Methods(http.MethodGet)
	return
}
