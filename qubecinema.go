package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type PersonResponse struct {
	Name   string        `json:"name"`
	Movies []UrlResponse `json:"movies"`
}

type MovieResponse struct {
	Name string        `json:"name"`
	Cast []UrlResponse `json:"cast"`
	Crew []UrlResponse `json:"crew"`
}

type UrlResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Role string `json:"role"`
}

type PersonNode struct {
	Name          string
	Identifier    string
	PreviousRole  string
	NextMovieRole string
	PreviousMovie *MovieNode
}

type MovieNode struct {
	Name           string
	Identifier     string
	PreviousPerson *PersonNode
}

func FindDegreeofSeparation(person1, person2 string) {
	err := validate(person1)
	if err != nil {
		fmt.Println("Invalid : " + person1)
		return
	}

	err = validate(person2)
	if err != nil {
		fmt.Println("Invalid : " + person2)
		return
	}

	visited := make(map[string]bool)

	personQueue := make(map[string]*MovieNode)
	personQueue[person1] = nil

	movieQueue := make(map[string]*PersonNode)

	degree := 0

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for len(personQueue) > 0 {
		if val, ok := personQueue[person2]; ok {
			generateResponse(degree, person2, val)
			return
		}
		for key, value := range personQueue {
			wg.Add(1)
			go func(identifier string, previousMovie *MovieNode, movieQueue map[string]*PersonNode, wg *sync.WaitGroup, mutex *sync.Mutex) {
				defer wg.Done()
				tempResponse, err := getPersonResponse(identifier)
				if err != nil {
					return
				}
				tempNode := PersonNode{
					Name:          tempResponse.Name,
					Identifier:    identifier,
					PreviousMovie: previousMovie,
				}

				mutex.Lock()
				defer mutex.Unlock()
				if _, ok := visited[identifier]; ok {
					return
				}
				for _, movie := range tempResponse.Movies {
					if previousMovie != nil && movie.Url == previousMovie.Identifier {
						tempNode.PreviousRole = movie.Role
					}
					if _, ok := visited[movie.Url]; !ok {
						movieQueue[movie.Url] = &tempNode
					}
				}
				visited[identifier] = true
			}(key, value, movieQueue, &wg, &mutex)
		}
		wg.Wait()
		personQueue = make(map[string]*MovieNode)
		time.Sleep(time.Duration(degree) * time.Second)

		if len(movieQueue) > 0 {
			for key, value := range movieQueue {
				wg.Add(1)
				go func(identifier string, previousPerson *PersonNode, personQueue map[string]*MovieNode, wg *sync.WaitGroup, mutex *sync.Mutex) {
					defer wg.Done()
					tempResponse, err := getMovieResponse(identifier)
					if err != nil {
						return
					}
					tempNode := MovieNode{
						Name:           tempResponse.Name,
						Identifier:     identifier,
						PreviousPerson: previousPerson,
					}

					mutex.Lock()
					defer mutex.Unlock()
					if _, ok := visited[identifier]; ok {
						return
					}
					for _, cast := range tempResponse.Cast {
						if cast.Url == tempNode.PreviousPerson.Identifier {
							previousPerson.NextMovieRole = cast.Role
						}
						if _, ok := visited[cast.Url]; !ok {
							personQueue[cast.Url] = &tempNode
						}
					}

					for _, crew := range tempResponse.Crew {
						if crew.Url == tempNode.PreviousPerson.Identifier {
							previousPerson.NextMovieRole = crew.Role
						}
						if _, ok := visited[crew.Url]; !ok {
							personQueue[crew.Url] = &tempNode
						}
					}

					visited[identifier] = true
				}(key, value, personQueue, &wg, &mutex)
			}
			wg.Wait()
			movieQueue = make(map[string]*PersonNode)
			time.Sleep(time.Duration(degree) * time.Second)
		}
		degree++
	}
}

func getPersonResponse(identifier string) (*PersonResponse, error) {
	url := "https://data.moviebuff.com/" + identifier
	resp, err := http.Get(url)
	if err != nil {
		resp, err = http.Get(url)
		if err != nil {
			return nil, errors.New("Couldn't get response for : " + identifier)
		}
	}
	var person PersonResponse
	json.NewDecoder(resp.Body).Decode(&person)
	return &person, nil
}

func getMovieResponse(identifier string) (*MovieResponse, error) {
	url := "https://data.moviebuff.com/" + identifier
	resp, err := http.Get(url)
	if err != nil {
		time.Sleep(1 * time.Second)
		resp, err = http.Get(url)
		if err != nil {
			return nil, errors.New("Couldn't get response for : " + identifier)
		}
	}
	var movie MovieResponse
	json.NewDecoder(resp.Body).Decode(&movie)
	return &movie, nil
}

func generateResponse(degree int, person2 string, movie *MovieNode) {
	fmt.Printf("Degree of Separation : %d\n", degree)
	var output []string
	movieResponse, err := getMovieResponse(movie.Identifier)
	if err != nil {
		fmt.Println("Error getting movie response : " + movie.Identifier)
	}

	var role string
	var name string

	for _, cast := range movieResponse.Cast {
		if cast.Url == person2 {
			role = cast.Role
			name = cast.Name
			break
		}
	}

	for _, crew := range movieResponse.Crew {
		if crew.Url == person2 {
			role = crew.Role
			name = crew.Name
			break
		}
	}

	for movie != nil {
		output = append(output, fmt.Sprintf("%s : %s", role, name))
		output = append(output, fmt.Sprintf("%s : %s", movie.PreviousPerson.NextMovieRole, movie.PreviousPerson.Name))
		output = append(output, fmt.Sprintf("Movie : %s", movie.Name))
		role = movie.PreviousPerson.PreviousRole
		name = movie.PreviousPerson.Name
		movie = movie.PreviousPerson.PreviousMovie
	}

	for i := len(output) - 1; i >= 0; i -= 3 {
		fmt.Println()
		fmt.Println(output[i])
		fmt.Println(output[i-1])
		fmt.Println(output[i-2])
	}
}

func validate(person string) error {
	_, err := getPersonResponse(person)
	return err
}
