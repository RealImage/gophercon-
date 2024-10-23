package main

import (
	"encoding/json"
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
		log.Println("Error fetching movie data: ", err)
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
