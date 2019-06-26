package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Get all users
/* func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
} */

// Get single user
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	processStart := time.Now()

	VPK := getPublicKey() // get Public Key
	params := mux.Vars(r) // Gets params
	userId := params["userId"]
	db := dbConn()
	defer db.Close()
	// ------------------------------------------------------------------------------
	// Get User

	row, err := db.Query("SELECT email, name, surnames FROM valentium.users WHERE userId = ?", userId)
	if err != nil {
		response := ErrorResponse{Status: "Error", Error: err.Error(), ExecutionTime: time.Since(processStart).Seconds() * 1000}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer row.Close()

	user := User{}
	count := 0

	var email, name, surnames []byte

	if row.Next() {
		err = row.Scan(&email, &name, &surnames)
		count += 1
		if err != nil {
			response := ErrorResponse{Status: "Error", Error: err.Error(), ExecutionTime: time.Since(processStart).Seconds() * 1000}
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	if count == 0 {
		response := ErrorResponse{Status: "Error", Error: "User not found", ExecutionTime: time.Since(processStart).Seconds() * 1000}
		json.NewEncoder(w).Encode(response)
		return
	}

	refId := createHash(createHash(userId) + createHash(string(VPK)))

	// ------------------------------------------------------------------------------
	// Get UPK

	upkSeed, err := getUserPrivateKey(refId)
	if err != nil {
		response := ErrorResponse{Status: "Error", Error: "Cannot get UPK", ExecutionTime: time.Since(processStart).Seconds() * 1000}
		json.NewEncoder(w).Encode(response)
		return
	}

	upk := createHash32(string(upkSeed) + createHash(string(VPK)))
	//log.Println("upk", string(thisUpk))

	// ------------------------------------------------------------------------------
	// Decrypt

	user.UserId = userId
	user.Email = decrypt(email, []byte(upk))
	user.Name = decrypt(name, []byte(upk))
	user.Surnames = decrypt(surnames, []byte(upk))

	//log.Println("user", user)
	//log.Println("Execution Time:", time.Since(processStart))

	response := UserResponse{Status: "OK", Data: user, ExecutionTime: time.Since(processStart).Seconds() * 1000}

	json.NewEncoder(w).Encode(response)

}

// Add new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Surnames string `json:"surnames"`
		Password string `json:"password"`
	}

	processStart := time.Now()

	VPK := getPublicKey() // get Public Key

	_ = json.NewDecoder(r.Body).Decode(&user)

	t := time.Now()

	userId := createHash(createHash(string(user.Email)+createHash(string(VPK))) + t.String())
	emailHash := createHash(createHash(string(user.Email) + createHash(string(VPK))))
	login := createHash(createHash(string(user.Email) + createHash(string(user.Password)) + createHash(string(VPK))))

	refId := createHash(createHash(userId) + createHash(string(VPK)))
	upkSeed := createHash32(createHash(refId) + createHash(string(VPK)))

	// ------------------------------------------------------------------------------
	// Save UPK

	_, err := saveUserPrivateKey(refId, upkSeed)
	if err != nil {
		response := ErrorResponse{Status: "Error", Error: "Cannot save UPK Seed", ExecutionTime: time.Since(processStart).Seconds() * 1000}
		json.NewEncoder(w).Encode(response)
		return
	}

	// ------------------------------------------------------------------------------
	// Encrypt

	upk := createHash32(upkSeed + createHash(string(VPK)))

	email := encrypt([]byte(user.Email), []byte(upk))
	name := encrypt([]byte(user.Name), []byte(upk))
	surnames := encrypt([]byte(user.Surnames), []byte(upk))

	// ------------------------------------------------------------------------------
	// Save User

	db := dbConn()

	_, err = db.Exec("INSERT INTO valentium.users (userId, email, emailHash, login, name, surnames, createdDate) VALUES (?,?,?,?,?,?,?)", userId, email, emailHash, login, name, surnames, t.String())
	if err != nil {
		response := ErrorResponse{Status: "OK", Error: err.Error(), ExecutionTime: time.Since(processStart).Seconds() * 1000}
		json.NewEncoder(w).Encode(response)
		return
	}
	defer db.Close()

	response := GenericResponse{Status: "OK", Data: map[string]string{"userId": userId}, ExecutionTime: time.Since(processStart).Seconds() * 1000}

	json.NewEncoder(w).Encode(response)

}

func saveUserPrivateKey(refId string, thisUpk string) (int, error) {

	/*
		// ------------------------------------------------------------------------------
		// Save UPK MySQL

			_, err := db.Exec("INSERT INTO valentium.keys (refId, upk) VALUES (?,?)", refId, thisUpk)
			if err != nil {
				response := ErrorResponse{Status: "OK", Error: err.Error(), ExecutionTime: time.Since(processStart).Seconds() * 1000}
				json.NewEncoder(w).Encode(response)
				return
			}

	*/

	// ------------------------------------------------------------------------------
	// Save UPK MongoDB

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBUri))
	if err != nil {
		return 0, err
	}
	collection := client.Database("valentium").Collection("wq")
	_, err = collection.InsertOne(ctx, bson.M{"r1": refId, "n1": thisUpk})
	if err != nil {
		return 0, err
	}
	//log.Println(result)
	defer client.Disconnect(ctx)
	return 0, nil

}

func getUserPrivateKey(refId string) ([]byte, error) {

	/*
		// ------------------------------------------------------------------------------
		// Get UPK from MySQL


			row1, err := db.Query("SELECT upk FROM valentium.keys WHERE refId = ?", refId)
			if err != nil {
				log.Fatal(err)
			}

			var thisUpk []byte
			count = 0

			if row1.Next() {
				row1.Scan(&thisUpk)
				count += 1
			}

			if count == 0 {
				response := ErrorResponse{Status: "Error", Error: "UPK not found", ExecutionTime: time.Since(processStart).Seconds() * 1000}
				json.NewEncoder(w).Encode(response)
				return
			}

			defer db.Close()
	*/

	// ------------------------------------------------------------------------------
	// Get UPK from MongoDB

	var upk UserPrivateKey

	// create a new context for MongoDB
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDBUri))
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)

	collection := client.Database("valentium").Collection("wq")
	err = collection.FindOne(ctx, bson.M{"r1": refId}).Decode(&upk)
	if err != nil {
		return nil, err
	}

	return []byte(upk.N1), nil

}

/*
// Update user
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			var user User
			_ = json.NewDecoder(r.Body).Decode(&user)
			user.ID = params["id"]
			users = append(users, user)
			json.NewEncoder(w).Encode(user)
			return
		}
	}
}

// Delete user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range users {
		if item.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}
*/
