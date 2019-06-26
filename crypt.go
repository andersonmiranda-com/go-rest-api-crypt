package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
)

/*
func example() {

	text := []byte("My name is Astaxie")
	key := []byte("the-key-has-to-be-32-bytes-long!")

	ciphertext, err := encrypt(text, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s => %x\n", text, ciphertext)

	plaintext, err := decrypt(ciphertext, key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%x => %s\n", ciphertext, plaintext)

}
*/

//32 bytes initial public key
var initKey []byte = []byte("640323df00dcb3aff871edb82381e8d7")

var userPrivateKey []byte

func getPublicKey() []byte {
	return initKey
}

func encrypt(plaintext []byte, key []byte) string {

	if len(plaintext) == 0 {
		return ""
	}

	if len(key) == 0 {
		return "Invalid key"
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return "Error encrypting"
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
		return "Error encrypting"
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "Error encrypting"
	}

	return string(gcm.Seal(nonce, nonce, plaintext, nil))
}

func decrypt(ciphertext []byte, key []byte) string {

	if len(ciphertext) == 0 {
		return ""
	}

	if len(key) == 0 {
		return "Invalid key"
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
		return "Error decrypting"
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		log.Fatal(err)
		return "Error decrypting"
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		if err != nil {
			log.Fatal("ciphertext too short")
			return "Error decrypting"
		}
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		log.Fatal(err)
		return "Error decrypting"
	}

	return string(plaintext)

}

/* func createHash32(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
*/
func createHash32(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	hashed := hex.EncodeToString(h.Sum(nil))
	return hashed[:32]
}

func createHash(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	hashed := hex.EncodeToString(h.Sum(nil))
	return hashed
}
