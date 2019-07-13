package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

//check email / pre-auth
func checkEmail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	processStart := time.Now()

	params := mux.Vars(r) // Gets params
	email := params["email"]
	user := User{}
	VPK := getPublicKey() // get Public Key
	emailHash := createHash(createHash(string(email) + createHash(string(VPK))))

	// ------------------------------------------------------------------------------
	// Get User
	db := dbConn()
	defer db.Close()

	row, err := db.Query("SELECT ou FROM valentium.users WHERE oh = ?", emailHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer row.Close()

	count := 0
	if row.Next() {
		err = row.Scan(&user.UserId)
		count += 1
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if count == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Email"))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		response := GenericResponse{Status: "OK", Data: map[string]string{"userId": user.UserId}, ExecutionTime: time.Since(processStart).Seconds() * 1000}
		json.NewEncoder(w).Encode(response)
	}

}

// Auth login/pass
func authenticate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	processStart := time.Now()
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)

	//log.Println(user.Email, user.Password)

	if len(user.Email) == 0 || len(user.Password) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please provide name and password"))
		return
	}

	VPK := getPublicKey() // get Public Key
	login := createHash(createHash(string(user.Email) + createHash(string(user.Password)) + createHash(string(VPK))))

	// ------------------------------------------------------------------------------
	// Get User

	db := dbConn()
	defer db.Close()

	row, err := db.Query("SELECT ou FROM valentium.users WHERE ca = ?", login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer row.Close()

	count := 0

	if row.Next() {
		err = row.Scan(&user.UserId)
		count += 1
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if count == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid login"))
		return
	} else {
		// ------------------------------------------------------------------------------
		// Generate JWT
		token, err := getToken(user.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error generating JWT token: " + err.Error()))
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)

			response := GenericResponse{Status: "OK", Data: map[string]string{"userId": user.UserId, "token": token}, ExecutionTime: time.Since(processStart).Seconds() * 1000}
			json.NewEncoder(w).Encode(response)
		}

	}

}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := verifyToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		userId := claims.(jwt.MapClaims)["userId"].(string)
		r.Header.Set("userId", userId)
		next(w, r)
	})
}
