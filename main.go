package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
)

type Person struct {
	Name   string        `json:"name"`
	URL    string        `json:"url"`
	Movies []PersonMovie `json:"movies"`
}

type PersonMovie struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type Movie struct {
	Name string        `json:"name"`
	URL  string        `json:"url"`
	Cast []MoviePerson `json:"cast"`
	Crew []MoviePerson `json:"crew"`
}

type MoviePerson struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

var cache = make(map[string]interface{})
var cacheLock = &sync.Mutex{}

func fetchData(url string) ([]byte, error) {
	resp, err := http.Get("https://data.moviebuff.com/" + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

type Path struct {
	Depth []Step
}

type Step struct {
	Movie  string
	Role   string
	Person string
}

func findConnection(start, end string) (*Path, error) {
	// Find Shortest Path from single source using BFS with Queue iteratively
	visited := make(map[string]bool)
	queue := make([]Path, 0)

	queue = append(queue, Path{Depth: []Step{}})
	visited[start] = true

	for len(queue) > 0 {
		currentPath := queue[0]
		queue = queue[1:]

		currentPerson := start
		if len(currentPath.Depth) > 0 {
			currentPerson = currentPath.Depth[len(currentPath.Depth)-1].Person
		}

		// Check if data exists in cache before making API call
		var person Person
		data, exists := cache[currentPerson]
		if !exists {
			rawData, err := fetchDataAsync(currentPerson)
			if err != nil {
				continue
			}
			if err := json.Unmarshal(rawData, &person); err != nil {
				continue
			}
			// Store the response in Cache for future use
			// Using Lock to make map thread-safe
			cacheLock.Lock()
			cache[currentPerson] = person
			cacheLock.Unlock()
		} else {
			person = data.(Person)
		}

		var wg sync.WaitGroup
		movieChan := make(chan Movie, len(person.Movies))

		for _, movie := range person.Movies {
			wg.Add(1)
			// Fetching movie data asynchronously by spawning goroutine
			go func(movie PersonMovie) {
				defer wg.Done()
				var movieData Movie
				data, exists := cache[movie.URL]
				if !exists {
					rawData, err := fetchDataAsync(movie.URL)
					if err != nil {
						return
					}
					if err := json.Unmarshal(rawData, &movieData); err != nil {
						return
					}
					cacheLock.Lock()
					cache[movie.URL] = movieData
					cacheLock.Unlock()
				} else {
					movieData = data.(Movie)
				}
				movieChan <- movieData
			}(movie)
		}

		go func() {
			wg.Wait()
			close(movieChan)
		}()

		for movieData := range movieChan {
			// Check in both cast and crew of the movie for connections
			for _, p := range append(movieData.Cast, movieData.Crew...) {
				if visited[p.URL] {
					continue
				}

				visited[p.URL] = true

				newPath := Path{
					Depth: append(currentPath.Depth, Step{
						Movie:  movieData.Name,
						Role:   p.Role,
						Person: p.URL,
					}),
				}

				if p.URL == end {
					return &newPath, nil
				}

				queue = append(queue, newPath)
			}
		}
	}

	return nil, fmt.Errorf("connection error")
}

func fetchDataAsync(url string) ([]byte, error) {
	dataChan := make(chan []byte, 1)
	errChan := make(chan error, 1)

	go func() {
		data, err := fetchData(url)
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- data
	}()

	select {
	case data := <-dataChan:
		return data, nil
	case err := <-errChan:
		return nil, err
	}
}

func main() {
	var person1, person2 string

	fmt.Print("source actor : ")
	fmt.Scanf("%s", &person1)
	fmt.Print("target actor : ")
	fmt.Scanf("%s", &person2)
	path, err := findConnection(person1, person2)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Printf("Degrees of Separation: %d\n\n", len(path.Depth))
	for i, step := range path.Depth {
		fmt.Printf("%d. Movie: %s\n", i+1, step.Movie)
		if i == 0 {
			fmt.Printf("   %s: %s\n", step.Role, person1)
		} else {
			fmt.Printf("   %s: %s\n", path.Depth[i-1].Role, path.Depth[i-1].Person)
		}
		fmt.Printf("   %s: %s\n", step.Role, step.Person)
	}
}
