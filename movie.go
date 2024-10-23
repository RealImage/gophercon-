package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type PersonRole struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type Movie struct {
	URL  string       `json:"url"`
	Type string       `json:"type"`
	Name string       `json:"name"`
	Cast []PersonRole `json:"cast"`
	Crew []PersonRole `json:"crew"`
}

func FetchMovieDetails(movieURL string) (Movie, error) {
	data, err := FetchData(movieURL)
	if err != nil {
		return Movie{}, err
	}

	var movie Movie
	err = json.Unmarshal(data, &movie)
	if err != nil {
		log.Println("Error unmarshalling movie data: ", err)
		return Movie{}, err
	}

	return movie, nil
}

func MoviesPeople(movieRoles []MovieRole) map[string]PersonRole {
	moviePeople := make(map[string]PersonRole)
	for _, movieRole := range movieRoles {
		movieData, err := FetchMovieDetails(movieRole.URL)
		if err != nil {
			continue
		}

		for _, person := range movieData.Cast {
			moviePeople[person.URL] = person
		}

		for _, person := range movieData.Crew {
			moviePeople[person.URL] = person
		}
	}

	return moviePeople
}

func HasSamePerson(moviesPeople1 map[string]PersonRole, moviesPeople2 map[string]PersonRole) bool {
	for personURL := range moviesPeople1 {
		if _, ok := moviesPeople2[personURL]; ok {
			fmt.Println(moviesPeople1[personURL].URL)
			return true
		}
	}

	for personURL := range moviesPeople2 {
		if _, ok := moviesPeople1[personURL]; ok {
			fmt.Println(moviesPeople1[personURL].URL)
			return true
		}
	}

	return false
}
