package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Person struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Hobby *Hobby `json:"hobby"`
}

type Hobby struct {
	Fav string `json:"fav"`
}

var people []Person

// jsonMiddleware sets Content-Type to application/json for all responses
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Sample data
	people = append(people, Person{Name: "Tarsh", ID: "69", Hobby: &Hobby{Fav: "badminton"}})

	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	peopleRouter := router.PathPrefix("/people").Subrouter()
	peopleRouter.HandleFunc("/", getPeople).Methods("GET")
	peopleRouter.HandleFunc("/{id}", getPersonById).Methods("GET")
	peopleRouter.HandleFunc("/", addPerson).Methods("POST")
	peopleRouter.HandleFunc("/{id}", editPersonById).Methods("PUT")
	peopleRouter.HandleFunc("/{id}", deletePersonById).Methods("DELETE")

	fmt.Println("Starting server on port 8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatal(err)
	}
}

func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

func getPersonById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	for _, item := range people {
		if item.ID == id {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Person not found", http.StatusNotFound)
}

func addPerson(w http.ResponseWriter, r *http.Request) {
	var newPerson Person
	_ = json.NewDecoder(r.Body).Decode(&newPerson)
	newId := rand.Intn(10000)
	newPerson.ID = strconv.Itoa(newId)
	people = append(people, newPerson)
	json.NewEncoder(w).Encode(people)
}

func editPersonById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	var personToEdit Person
	if err := json.NewDecoder(r.Body).Decode(&personToEdit); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Iterate to find the person and update it
	for index, item := range people {
		if item.ID == id {
			// Preserve the original ID
			personToEdit.ID = id
			// Update the person at the current index
			people[index] = personToEdit
			json.NewEncoder(w).Encode(personToEdit)
			return
		}
	}
	// If person with the given ID is not found, return an error
	http.Error(w, "Person not found", http.StatusNotFound)
}

func deletePersonById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	for index, person := range people {
		if person.ID == id {
			// Remove the person from the slice
			people = append(people[:index], people[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(people)
}
