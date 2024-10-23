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

func ToMovieRoleSet(movieRoles []MovieRole) map[string]MovieRole {
	movieRoleSet := make(map[string]MovieRole)
	for _, movieRole := range movieRoles {
		movieRoleSet[movieRole.URL] = movieRole
	}

	return movieRoleSet
}

func HasSameMovie(movieRoleSet1 map[string]MovieRole, movieRoleSet2 map[string]MovieRole) bool {
	for movieURL := range movieRoleSet1 {
		if _, ok := movieRoleSet2[movieURL]; ok {
			return true
		}
	}

	for movieURL := range movieRoleSet2 {
		if _, ok := movieRoleSet1[movieURL]; ok {
			return true
		}
	}

	return false
}
