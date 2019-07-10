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
	log.Println("Starting Valetium API server on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", r))

}

func ConfigureRouter() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/login", authenticate).Methods("POST")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", authMiddleware(getUser)).Methods("GET")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	//r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	//r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	return r

}
