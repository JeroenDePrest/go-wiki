package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var templates = template.Must(template.ParseFiles("tmpl/index.html"))
var session *mgo.Session

const authKey = "mysupersecretkey"

type Page struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

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
		"foo": "bar",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(authKey))

	return tokenString, err
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", nil)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t User
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = register(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var t User
	err := decoder.Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	token, err := auth(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + token + `"}`))
}

func pagesHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		if title := strings.Split(r.URL.Path, "/"); len(title) < 0 {
			p, err := loadPage(title[3])
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
		} else {
			p, err := loadAllPages()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
		}

	case "POST":
		decoder := json.NewDecoder(r.Body)
		var t Page
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		err = createPage(t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)

	}

}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func authMiddleware(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if bearer := r.Header.Get("Authorization"); bearer != "" {
			tokenString := strings.Split(bearer, " ")[1]
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(authKey), nil
			})

			if token.Valid {
				fn(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

	}
}

func main() {
	var err error
	session, err = mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	//To be able to server static files
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))

	//handling routes
	http.HandleFunc("/api/pages", authMiddleware(pagesHandler))
	//for a single page
	http.HandleFunc("/api/pages/", authMiddleware(pagesHandler))
	http.HandleFunc("/api/users", userHandler)
	http.HandleFunc("/api/auth", authHandler)
	http.HandleFunc("/", rootHandler)

	http.ListenAndServe(":8080", nil)
}
