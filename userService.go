package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"os"
	"time"
)

func auth(user User) (string, error) {
	result := User{}

	hash := sha256.Sum256([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])

	c := session.DB("godb").C("users")
	err := c.Find(bson.M{"name": user.Name, "password": user.Password}).One(&result)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("authKey")))

	return tokenString, err
}
