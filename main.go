//package main
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//Person  is an indivvidual person.
type Person struct {
	ID        string   `json:"id,omitempty"`
	FirstName string   `json:"firstrname,omitempty"`
	LastName  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

//Address is a address of person
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person

//GetPersonEndpoint  get a special person
func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _,item := range people{
		if item.ID == params["id"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Person{})
}

//GetPeopleEndpoint  get a slice of person
func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

//CreatePersonEndpoint create a person
func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people,person)
	json.NewEncoder(w).Encode(people)
}

//DeletePersonEndpoint  delete a person by id
func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range  people{
		if item.ID == params["id"]{
			people = append(people[:index],people[index+1:]...)
			json.NewEncoder(w).Encode(people)
			return
		}
	}
	json.NewEncoder(w).Encode(people)

}

func main() {
	router := mux.NewRouter()

	people = append(people, Person{ID: "1", FirstName: "Nic", LastName: "Raboy", Address: &Address{City: "Dublin", State: "California"}})
	people = append(people, Person{ID: "2", FirstName: "Jack", LastName: "Zhong", Address: &Address{City: "Dublin", State: "California"}})
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":12345", router))
}
