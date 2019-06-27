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
	r := mux.NewRouter()

	// Route handles & endpoints
	//r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{userId}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users/login", loginUser).Methods("POST")
	//r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	//r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":5000", r))

}
