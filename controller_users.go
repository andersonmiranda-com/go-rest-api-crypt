package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Get all users
/* func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
} */

// Add new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user User

	processStart := time.Now()

	VPK := getPublicKey() // get Public Key

	_ = json.NewDecoder(r.Body).Decode(&user)

	t := time.Now()

	userId := createHash(createHash(string(user.Email)+createHash(string(VPK))) + t.String())
	emailHash := createHash(createHash(string(user.Email) + createHash(string(VPK))))
	login := createHash(createHash(string(user.Email) + createHash(string(user.Password)) + createHash(string(VPK))))

	refId := createHash(createHash(userId) + createHash(string(VPK)))
	upkSeed := createHash(createHash(refId) + createHash(string(VPK)))

	// ------------------------------------------------------------------------------
	// Save UPK

	_, err := saveUserPrivateKeySeed(refId, upkSeed)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Cannot save UPK Seed"))
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

	_, err = db.Exec("INSERT INTO valentium.users (ou, oe, oh, ca, rn, bn, dc) VALUES (?,?,?,?,?,?,?)", userId, email, emailHash, login, name, surnames, t.String())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer db.Close()

	response := GenericResponse{Status: "OK", Data: map[string]string{"userId": userId}, ExecutionTime: time.Since(processStart).Seconds() * 1000}

	json.NewEncoder(w).Encode(response)

}

// Get single user
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	processStart := time.Now()

	//params := mux.Vars(r) // Gets params
	//userId := params["userId"]

	// ------------------------------------------------------------------------------
	// Get token and compares
	//tokenUserId := r.Header.Get("userId")
	//if userId != tokenUserId {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	w.Write([]byte("You dot have permission to access this data!"))
	//	return
	//}

	userId := r.Header.Get("userId")

	VPK := getPublicKey() // get Public Key
	db := dbConn()
	defer db.Close()

	// ------------------------------------------------------------------------------
	// Get User

	row, err := db.Query("SELECT oe, rn, bn, dc FROM valentium.users WHERE ou = ?", userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer row.Close()

	user := User{}
	count := 0

	var email, name, surnames []byte
	var createdDate string

	if row.Next() {
		err = row.Scan(&email, &name, &surnames, &createdDate)
		count += 1
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if count == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found"))
		return
	}

	refId := createHash(createHash(userId) + createHash(string(VPK)))

	// ------------------------------------------------------------------------------
	// Get UPK

	upkSeed, err := getUserPrivateKeySeed(refId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Cannot get UPK"))
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
	user.CreatedDate = createdDate
	//log.Println("user", user)
	//log.Println("Execution Time:", time.Since(processStart))

	response := UserResponse{Status: "OK", Data: user, ExecutionTime: time.Since(processStart).Seconds() * 1000}

	json.NewEncoder(w).Encode(response)

}

func saveUserPrivateKeySeed(refId string, thisUpk string) (int, error) {

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

func getUserPrivateKeySeed(refId string) ([]byte, error) {

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
