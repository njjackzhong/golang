//package main
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/elgs/jsonql"
	"fmt"
)

var jsonString = `[{"name": "elgs"},{"name": "enny"},{"name": "sam"}]`
//var jsonString = `
//[
//  {
//    "name": "elgs",
//    "gender": "m",
//	"age": 35,
//    "skills": [
//      "Golang",
//      "Java",
//      "C"
//    ]
//  },
//  {
//    "name": "enny",
//    "gender": "f",
//    "age": 36,
//	"skills": [
//      "IC",
//      "Electric design",
//      "Verification"
//    ]
//  },
//  {
//    "name": "sam",
//    "gender": "m",
//	"age": 1,
//    "skills": [
//      "Eating",
//      "Sleeping",
//      "Crawling"
//    ]
//  }
//]
//`
type Message struct {
	code    string `json:"code,omitempty"`
	message string `json:"message,omitempty"`
}

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
	for _, item := range people {
		if item.ID == params["id"] {
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
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

//DeletePersonEndpoint  delete a person by id
func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			json.NewEncoder(w).Encode(people)
			return
		}
	}
	json.NewEncoder(w).Encode(people)

}

func SelectMessage(w http.ResponseWriter, req *http.Request) {

	messages := req.URL.Query().Get("messages")
	where := req.URL.Query().Get("where")

	fmt.Println("messages: " + messages)
	fmt.Println("where: " + where)

	parser, err := jsonql.NewStringQuery(messages)

	if err == nil {
		selectedMessages, err := parser.Query(where)

		fmt.Println(selectedMessages)
		fmt.Println(err)
		//if err == nil {
		//json.NewEncoder(w).Encode(true)

		json.NewEncoder(w).Encode(selectedMessages)
		//}

	}

}

func main() {
	//jsonQlTest()

	router := mux.NewRouter()

	people = append(people, Person{ID: "1", FirstName: "Nic", LastName: "Raboy", Address: &Address{City: "Dublin", State: "California"}})
	people = append(people, Person{ID: "2", FirstName: "Jack", LastName: "Zhong", Address: &Address{City: "Dublin", State: "California"}})
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	router.HandleFunc("/select", SelectMessage).Methods("GET")

	log.Fatal(http.ListenAndServe(":12345", router))
}

func jsonQlTest() {
	parser, err := jsonql.NewStringQuery(jsonString)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(parser.Query("name='elgs'"))
	//[map[skills:[Golang Java C] name:elgs gender:m age:35]] <nil>

	//fmt.Println(parser.Query("name='elgs' && gender='f'"))
	////[] <nil>
	//
	//fmt.Println(parser.Query("age<10 || (name='enny' && gender='f')"))
	////[map[gender:f age:36 skills:[IC Electric design Verification] name:enny] map[age:1 skills:[Eating Sleeping Crawling] name:sam gender:m]] <nil>
	//
	//fmt.Println(parser.Query("age<10"))
	////[map[name:sam gender:m age:1 skills:[Eating Sleeping Crawling]]] <nil>
	//
	//fmt.Println(parser.Query("1=0"))
	////[] <nil>
	//
	//fmt.Println(parser.Query("age=(2*3)^2"))
	////[map[skills:[IC Electric design Verification] name:enny gender:f age:36]] <nil>
	//
	//fmt.Println(parser.Query("name ~= 'e.*'"))
	////[map[age:35 skills:[Golang Java C] name:elgs gender:m] map[skills:[IC Electric design Verification] name:enny gender:f age:36]] <nil>
	//
	//fmt.Println(parser.Query("name='el'+'gs'"))
	//fmt.Println(parser.Query("age=30+5.0"))
	//fmt.Println(parser.Query("age=40.0-5"))
	//fmt.Println(parser.Query("age=70-5*7"))
	//fmt.Println(parser.Query("age=70.0/2.0"))
	//fmt.Println(parser.Query("age=71%36"))
}
