package page

import (
	"net/http"
	"github.com/gorilla/mux"
	"encoding/json"
)

func PagesHandler(w http.ResponseWriter, r *http.Request) {

	p, err := loadAllPages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)

}

func PageHandler(w http.ResponseWriter, r *http.Request) {

	title := mux.Vars(r)["title"]
	p, err := loadPage(title)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)

}

func CreatePageHandler(w http.ResponseWriter, r *http.Request) {

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
