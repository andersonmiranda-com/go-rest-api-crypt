package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

type UserTest struct {
	Name     string `json:"name"`
	Surnames string `json:"surnames"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Main function
func TestRegisterUsers(t *testing.T) {

	processStart := time.Now()

	// Open our jsonFile
	jsonFile, err := os.Open("users-10k.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users-10k.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var users []UserTest

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &users)

	// we iterate through every user within our users array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for i := 0; i < len(users); i++ {
		jsonValue, _ := json.Marshal(users[i])
		resp, err := http.Post("http://api.valentium.io/users", "application/json", bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Println(err)
		}
		log.Println(i, users[i])
		body, err := ioutil.ReadAll(resp.Body)
		log.Println(string(body))
		//time.Sleep(250 * time.Millisecond)
	}

	log.Println("Total Time (s):", time.Since(processStart).Seconds())
}
