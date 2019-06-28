package main

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func getToken(userId string) (string, error) {
	VPK := getPublicKey()

	expirationTime := time.Now().Add(5 * time.Minute).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp":    expirationTime})

	tokenString, err := token.SignedString(VPK)
	return tokenString, err
}

func verifyToken(tokenString string) (jwt.Claims, error) {
	VPK := getPublicKey()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return VPK, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
