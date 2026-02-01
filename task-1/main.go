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
	router := http.NewServeMux()

	router.HandleFunc("GET /", helloHandler)
	router.HandleFunc("GET /categories", getCategoriesHandler)
	router.HandleFunc("GET /categories/{id}", getCategoryHandler)

	router.HandleFunc("POST /categories", postCategoriesHandler)
	router.HandleFunc("PUT /categories/{id}", putCategoryHandler)
	router.HandleFunc("DELETE /categories/{id}", deleteCategoryHandler)

	fmt.Println("Starting")
	if err := http.ListenAndServe(":8080", router); err != nil {
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get last id
	lastID := categories[len(categories)-1].ID

	w.Header().Set("Content-Type", "application/json")

	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	category.ID = lastID + 1
	categories = append(categories, category)

	json.NewEncoder(w).Encode(category)
}

func putCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
			category.ID = id
			categories[i] = category
			json.NewEncoder(w).Encode("Category updated")
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
			json.NewEncoder(w).Encode("Category deleted")
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
