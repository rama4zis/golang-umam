package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{
		ID:          1,
		Name:        "Category 1",
		Description: "Description 1",
	},
	{
		ID:          2,
		Name:        "Category 2",
		Description: "Description 2",
	},
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/categories", getCategoriesHandler)
	http.HandleFunc("/categories/", getCategoryHandler)

	http.HandleFunc("/categories", postCategoriesHandler)
	http.HandleFunc("/categories/", putCategoryHandler)
	http.HandleFunc("/categories/", deleteCategoryHandler)

	fmt.Println("Starting")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprint(w, "Hello World")
}

func getCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(categories)
}

func getCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, category := range categories {
		if category.ID == id {
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func postCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	categories = append(categories, category)

	json.NewEncoder(w).Encode(category)
}

func putCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories[i] = category
			json.NewEncoder(w).Encode(category)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
