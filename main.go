package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

type Fruit struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Color  string `json:"color,omitempty"`
	Rating int    `json:"rating,omitempty"`
}

var fruits []Fruit

func serveIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "API")
}

func GetFruits(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	err := enc.Encode(fruits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetFruit(w http.ResponseWriter, r *http.Request) {
	fruitID := chi.URLParam(r, "fruitID")
	for _, item := range fruits {
		if item.ID == fruitID {
			err := json.NewEncoder(w).Encode(item)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
	http.NotFound(w, r)
}

func AddFruit(w http.ResponseWriter, r *http.Request) {
	fruit := Fruit{}
	if err := json.NewDecoder(r.Body).Decode(&fruit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fruits = append(fruits, fruit)
	json.NewEncoder(w).Encode(fruits)
}

func DeleteFruit(w http.ResponseWriter, r *http.Request) {
	fruitID := chi.URLParam(r, "fruitID")
	for index, item := range fruits {
		if item.ID == fruitID {
			fruits = append(fruits[:index], fruits[index+1:]...)
			json.NewEncoder(w).Encode(fruits)
			return
		}
	}
	http.NotFound(w, r)
}

func UpdateFruit(w http.ResponseWriter, r *http.Request) {
	fruitID := chi.URLParam(r, "fruitID")
	for index, item := range fruits {
		if item.ID == fruitID {
			fruit := Fruit{}
			if err := json.NewDecoder(r.Body).Decode(&fruit); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fruits[index] = fruit
			json.NewEncoder(w).Encode(fruit)
			return
		}
	}
	http.NotFound(w, r)
}

func main() {
	fruits = append(fruits, Fruit{ID: "1", Name: "Orange", Color: "Orange", Rating: 1})
	fruits = append(fruits, Fruit{ID: "2", Name: "Banana", Color: "Yellow", Rating: 5})

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", serveIndex)
	r.Get("/fruits", GetFruits)
	r.Get("/fruits/:fruitID", GetFruit)
	r.Post("/fruits", AddFruit)
	r.Delete("/fruits/:fruitID", DeleteFruit)
	r.Put("/fruits/:fruitID", UpdateFruit)
	http.ListenAndServe(":7000", r)
}
