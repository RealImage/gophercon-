package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CommonData is the common data that exists in both Person and Movie.
type CommonData struct {
	URL  string `json:"url"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// Participation is the data that represents both:
// 1. A movie in which a person has played a role.
// 2. A person who played a role in a movie.
type Participation struct {
	URL  string `json:"url"`
	Name string `json:"name"`
	Role string `json:"role"`
}

// Person is the data that represents a person.
type Person struct {
	CommonData
	Movies []Participation `json:"movies"`
}

// Movie is the data that represents a movie.
type Movie struct {
	CommonData
	Cast []Participation `json:"cast"`
	Crew []Participation `json:"crew"`
}

// baseURL is the base URL for the Moviebuff API.
const baseURL = "https://data.moviebuff.com"

// FetchData is a helper function that fetches data from any Moviebuff URL.
func FetchData(moviebuffURL string) ([]byte, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", baseURL, moviebuffURL))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 200:
		return data, nil

	default:
		return nil, fmt.Errorf("invalid response. Status code: %d. Data: %s", res.StatusCode, string(data))
	}
}

// FetchMovie is a helper function that fetches movie data from a Moviebuff URL.
func FetchMovie(movieURL string) (Movie, error) {
	data, err := FetchData(movieURL)
	if err != nil {
		return Movie{}, err
	}

	var movie Movie
	err = json.Unmarshal(data, &movie)
	if err != nil {
		return Movie{}, err
	}

	return movie, nil
}

// FetchPerson is a helper function that fetches person data from a Moviebuff URL.
func FetchPerson(personURL string) (Person, error) {
	data, err := FetchData(personURL)
	if err != nil {
		return Person{}, err
	}

	var person Person
	err = json.Unmarshal(data, &person)
	if err != nil {
		return Person{}, err
	}

	return person, nil
}
