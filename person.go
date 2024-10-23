package main

import (
	"encoding/json"
	"log"
)

type MovieRole struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type Person struct {
	URL    string      `json:"url"`
	Type   string      `json:"type"`
	Name   string      `json:"name"`
	Movies []MovieRole `json:"movies"`
}

func FetchPersonDetails(personURL string) (Person, error) {
	data, err := FetchData(personURL)
	if err != nil {
		log.Println("Error fetching person data: ", err)
		return Person{}, err
	}

	var person Person
	err = json.Unmarshal(data, &person)
	if err != nil {
		log.Println("Error unmarshalling person data: ", err)
		return Person{}, err
	}

	return person, nil
}
