package main

import (
	"crypto/sha256"
	"encoding/hex"
	"gopkg.in/mgo.v2/bson"
)

func loadPage(title string) (Page, error) {
	result := Page{}
	c := session.DB("godb").C("pages")
	err := c.Find(bson.M{"title": title}).One(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

func createPage(page Page) error {
	c := session.DB("godb").C("pages")
	err := c.Insert(page)
	if err != nil {
		return err
	}
	return nil
}

func register(user User) error {
	hash := sha256.Sum256([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	c := session.DB("godb").C("users")
	err := c.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func loadAllPages() ([]Page, error) {
	result := make([]Page, 1)
	c := session.DB("godb").C("pages")
	err := c.Find(nil).All(&result)
	if err != nil {
		return result, err
	}
	return result, err
}
