package main

import (
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

var session *mgo.Session

func main() {
	//env variable
	os.Setenv("authKey", "mysupersecretkey")
	var err error
	session, err = mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	/*p := r.PathPrefix("/api/pages").Subrouter()
	p.HandleFunc("", authMiddleware(pagesHandler)).Methods("POST")
	p.HandleFunc("", authMiddleware(createPageHandler)).Methods("GET")
	p.HandleFunc("/{title}", authMiddleware(pageHandler)).Methods("GET")

	r.HandleFunc("/api/users", userHandler).Methods("POST")
	r.HandleFunc("/api/auth", authHandler).Methods("POST")
	http.Handle("/", r)*/

	router := CreateRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
