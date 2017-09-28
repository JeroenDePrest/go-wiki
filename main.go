package main

import (
	"gowiki/mongoDB"
	"gowiki/router"
	"log"
	"net/http"
)

func main() {

	mongoDB.Create()

	defer mongoDB.Session().Close()
	/*p := r.PathPrefix("/api/pages").Subrouter()
	p.HandleFunc("", authMiddleware(pagesHandler)).Methods("POST")
	p.HandleFunc("", authMiddleware(createPageHandler)).Methods("GET")
	p.HandleFunc("/{title}", authMiddleware(pageHandler)).Methods("GET")

	r.HandleFunc("/api/users", userHandler).Methods("POST")
	r.HandleFunc("/api/auth", authHandler).Methods("POST")
	http.Handle("/", r)*/

	router := router.Create()

	log.Fatal(http.ListenAndServe(":8080", router))
}
