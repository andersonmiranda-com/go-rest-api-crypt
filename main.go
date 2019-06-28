package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var VPK string

// Main function
func main() {

	// Init router
	r := ConfigureRouter()

	// Start server
	log.Fatal(http.ListenAndServe(":5000", r))

}

func ConfigureRouter() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/login", authenticate).Methods("POST")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/{userId}", authMiddleware(getUser)).Methods("GET")
	//r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	//r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	return r

}
