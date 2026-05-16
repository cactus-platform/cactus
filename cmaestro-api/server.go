package main

import (
	"cmaestro-api/internal/router"
	cingest "cmaestro-ingest"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := router.New()

	// Middlewares Configuration
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/users", getUsers)

	// **************************************** SERVICES ****************************************
	// system
	// health
	// status

	// registry
	// List every image hosted on local docker registry

	// registry/create (update-in-place)
	// registry/list
	// registry/delete
	// registry/_/function_path

	// insight/
	// insight/_/function_path
	// ******************************************************************************************
	// All these previous methods could be authenticated and authorised (feature available later)
	// ******************************************************************************************

	//r.ListenAndServe(":8080")

	cactuizedFunctions := cingest.Ingest(`from cactuskit import ApiMethod, ApiProtocol, HttpStatus, cactuize

@cactuize()
def simple_entrypoint():
    return f"Hello World from {simple_entrypoint}"`)

	fmt.Println(cactuizedFunctions)
}

func getUsers(w http.ResponseWriter, _ *http.Request) {
	users := []string{"admin"}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
