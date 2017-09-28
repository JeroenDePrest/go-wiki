package page

import (
	"gopkg.in/mgo.v2/bson"
	"gowiki/mongoDB"
)

func loadPage(title string) (Page, error) {
	session := mongoDB.Session()
	result := Page{}
	c := session.DB("godb").C("pages")
	err := c.Find(bson.M{"title": title}).One(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

func createPage(page Page) error {
	session := mongoDB.Session()
	c := session.DB("godb").C("pages")
	err := c.Insert(page)
	if err != nil {
		return err
	}
	return nil
}

func loadAllPages() ([]Page, error) {
	result := make([]Page, 1)
	session := mongoDB.Session()
	c := session.DB("godb").C("pages")
	err := c.Find(nil).All(&result)
	if err != nil {
		return result, err
	}
	return result, err
}
