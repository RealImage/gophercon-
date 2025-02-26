package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qube_assignment/models"
	"qube_assignment/utils"
	"sync"
)

var (
	personCache = make(map[string]*models.Person)
	movieCache  = make(map[string]*models.Movie)
	cacheLock   sync.Mutex
)

// FetchData gets JSON data from the Moviebuff API
func FetchData(url string, result interface{}) error {
	err := utils.Limiter.Wait(context.Background()) // Apply rate limiting
	if err != nil {
		return fmt.Errorf("rate limit error: %v", err)
	}

	apiURL := fmt.Sprintf("https://data.moviebuff.com/%s", url)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add a User-Agent header to avoid request blocking
	req.Header.Set("User-Agent", "Go-Moviebuff-Client/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return fmt.Errorf("API error: 403 Forbidden - Request blocked by server")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}

func GetPersonData(personID string) (*models.Person, error) {
	cacheLock.Lock()
	if data, exists := personCache[personID]; exists {
		cacheLock.Unlock()
		return data, nil
	}
	cacheLock.Unlock()

	var person models.Person
	err := FetchData(personID, &person)
	if err != nil {
		return nil, err
	}

	cacheLock.Lock()
	personCache[personID] = &person
	cacheLock.Unlock()

	return &person, nil
}

func GetMovieData(movieID string) (*models.Movie, error) {
	cacheLock.Lock()
	if data, exists := movieCache[movieID]; exists {
		cacheLock.Unlock()
		return data, nil
	}
	cacheLock.Unlock()

	var movie models.Movie
	err := FetchData(movieID, &movie)
	if err != nil {
		return nil, err
	}

	cacheLock.Lock()
	movieCache[movieID] = &movie
	cacheLock.Unlock()

	return &movie, nil
}
