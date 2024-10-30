package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Models

type Actor struct {
	URL string
}

type MovieResponse struct {
	Name string `json:"name"`
	Crew []Actor `json:"crew"`
	Cast []Actor `json:"cast"`
}

type ActorResponse struct {
	Movies []Movie `json:"movies"`
}

type Movie struct {
	URL string `json:"url"`
}

type RelationNode struct {
	FirstPerson   Actor
	SecondPerson Actor
	Movie        string
}

type QueueNode struct {
	Person Actor
	Path   []RelationNode
	IsDest bool
}



type CallVisitedMutex struct {
	mu      sync.RWMutex
	visited map[string]bool
}

func NewCallVisitedMutex() *CallVisitedMutex {
	return &CallVisitedMutex{
		visited: make(map[string]bool),
	}
}

func (c *CallVisitedMutex) Add(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.visited[name] = true
}

func (c *CallVisitedMutex) Contains(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.visited[name]
}

// Moviebuff Client

type Client interface {
	FetchActor(actor string) (*ActorResponse, error)
	FetchMovie(movie string) (*MovieResponse, error)
}

type client struct {
	httpClient *http.Client
	logger     *log.Logger
}

func NewClient(logger *log.Logger) Client {
	return &client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: logger,
	}
}

func (c *client) makeHttpReq(suffix string) (*http.Response, error) {
	httpUrl := fmt.Sprintf("https://data.moviebuff.com/%s", suffix)
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		c.logger.Printf("Error creating request for %s: %v", suffix, err)
		return nil, err
	}
	req.Header.Set("User-Agent", "MyMovieApp/1.0") // Set a custom User-Agent
	return c.httpClient.Do(req)
}

func (c *client) FetchActor(actor string) (*ActorResponse, error) {
	res, err := c.makeHttpReq(actor)
	if err != nil {
		c.logger.Printf("Error making request for actor %s: %v", actor, err)
		return nil, err
	}
	defer res.Body.Close()



	var data ActorResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		c.logger.Printf("Error decoding response for actor %s: %v", actor, err)
		return nil, err
	}
	return &data, nil
}

func (c *client) FetchMovie(movie string) (*MovieResponse, error) {
	res, err := c.makeHttpReq(movie)
	if err != nil {
		c.logger.Printf("Error making request for movie %s: %v", movie, err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.logger.Printf("Error: received status code %d for movie %s", res.StatusCode, movie)
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var data MovieResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		c.logger.Printf("Error decoding response for movie %s: %v", movie, err)
		return nil, err
	}
	return &data, nil
}

// Separation Degree Service

type SeparationDegree interface {
	Find(source, dest Actor)
}

type separationDegreeService struct {
	client Client
	logger *log.Logger
}

func NewSeparationDegreeService(logger *log.Logger) SeparationDegree {
	return &separationDegreeService{
		client: NewClient(logger),
		logger: logger,
	}
}

func (s *separationDegreeService) Find(source, dest Actor) {
	visited := NewCallVisitedMutex()
	queue := []QueueNode{{Person: source}}
	degree := 1

	for len(queue) > 0 {
		nxtLvlQueue := []QueueNode{}

		for _, curNode := range queue {
			if visited.Contains(curNode.Person.URL) {
				continue
			}
			visited.Add(curNode.Person.URL)

			resp, err := s.client.FetchActor(curNode.Person.URL)
			if err != nil {
				s.logger.Printf("Error fetching actor details for %s: %v", curNode.Person.URL, err)
				continue
			}

			for _, movie := range resp.Movies {
				movieDetails, err := s.client.FetchMovie(movie.URL)
				if err != nil {
					s.logger.Printf("Error fetching movie details for %s: %v", movie.URL, err)
					continue
				}

				for _, actor := range movieDetails.Cast {
					relationNode := RelationNode{
						FirstPerson:   curNode.Person,
						SecondPerson: actor,
						Movie:        movieDetails.Name,
					}

					if actor.URL == dest.URL {
						s.logger.Printf("Found degree of separation: %d between %s and %s through %s",
							degree, source.URL, dest.URL, movieDetails.Name)
						s.printRelation(relationNode, degree, curNode.Path)
						return
					}

					nxtLvlQueue = append(nxtLvlQueue, QueueNode{
						Person: actor,
						Path:   append(curNode.Path, relationNode),
					})
				}
			}
		}
		queue = nxtLvlQueue
		s.logger.Printf("Degree of separation: %d for %s and %s", degree, source.URL, dest.URL)
		degree++
	}

	s.logger.Println("No relation found.")
}

func (s *separationDegreeService) printRelation(relation RelationNode, degree int, path []RelationNode) {
	fmt.Printf("Degree of Separation: %d\n", degree)

	// Iterate through the path and print each relation
	for i, node := range append(path, relation) {
		if i > 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%d. Movie: %s\n", i+1, node.Movie)
		if i == 0 {
			fmt.Printf("Actor: %s\n", node.FirstPerson.URL)
			fmt.Printf("Supporting Actor: %s\n", node.SecondPerson.URL)
		} else {
			fmt.Printf("Supporting Actor: %s\n", node.FirstPerson.URL)
			fmt.Printf("Actor: %s\n", node.SecondPerson.URL)
		}
	}
}

// Main Function

func main() {
	// Set up logging to a .txt file
	logFile, err := os.OpenFile("app.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.Ldate|log.Ltime)

	if len(os.Args) < 3 {
		logger.Fatal("Usage: degrees <source_actor_url> <dest_actor_url>")
	}

	sourceURL := os.Args[1]
	destURL := os.Args[2]

	if sourceURL == destURL {
		logger.Fatal("Error: Source and destination URLs must be different.")
	}

	svc := NewSeparationDegreeService(logger)
	svc.Find(Actor{URL: sourceURL}, Actor{URL: destURL})
}
