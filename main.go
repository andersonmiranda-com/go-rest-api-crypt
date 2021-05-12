package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var VPK string

// Main function
func main() {

	// Init router
	r := ConfigureRouter()

	//headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// Start server
	log.Println("Starting Valetium API server on port 5000...")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))

}

func ConfigureRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/auth/email/{email}", checkEmail).Methods("GET")
	r.HandleFunc("/auth/login", authenticate).Methods("POST")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", authMiddleware(getUser)).Methods("GET")
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	//r.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	//r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	return r

}
