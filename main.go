package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

var actorCache = make(map[string]Actor)
var movieCache = make(map[string]Movie)

type Node struct {
	Name string
	Path []string
}

type Movie struct {
	MovieName string  `json:"name"`
	MovieURL  string  `json:"url"`
	MovieRole string  `json:"role"`
	MovieCast []Actor `json:"cast"`
	MovieCrew []Actor `json:"crew"`
}

type Actor struct {
	ActorURL string  `json:"url"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Role     string  `json:"role"`
	Movies   []Movie `json:"movies"`
}

func findDegreesofActor(actor1, actor2 string) ([]string, []string, error) {
	role := []string{"Actor"}
	queue := []Node{{Name: actor1, Path: []string{actor1}}}
	visited := make(map[string]bool)
	visited[actor1] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		person, err := getactordetails(current.Name)
		if err != nil {
			continue
		}

		for _, movie := range person.Movies {
			movieData, err := getMoviedetails(movie.MovieURL)

			if err != nil {
				continue
			}

			for _, cast := range append(movieData.MovieCast, movieData.MovieCrew...) {

				if cast.ActorURL == actor2 {
					return append(current.Path, movieData.MovieName, actor2), append(role, cast.Role), nil
				}

			}

			for _, actor := range append(movieData.MovieCast, movieData.MovieCrew...) {

				if !visited[actor.ActorURL] {
					visited[actor.ActorURL] = true
					queue = append(queue, Node{Name: actor.ActorURL, Path: append(current.Path, movieData.MovieName, actor.ActorURL)})
				}

			}

		}
	}

	return nil, []string{}, errors.New("no connection found")
}

func getMoviedetails(moviebuffURL string) (Movie, error) {
	var moviedetails Movie
	if movie, exists := movieCache[moviebuffURL]; exists {
		return movie, nil
	}

	url := fmt.Sprintf("https://data.moviebuff.com/%s", moviebuffURL)

	if err := getData(url, &moviedetails); err != nil {
		return Movie{}, err
	}

	movieCache[moviebuffURL] = moviedetails
	return moviedetails, nil
}

func getactordetails(actorName string) (Actor, error) {
	var actordetails Actor
	if person, exists := actorCache[actorName]; exists {
		return person, nil
	}

	url := fmt.Sprintf("https://data.moviebuff.com/%s", actorName)
	if err := getData(url, &actordetails); err != nil {
		return Actor{}, err
	}
	actorCache[actorName] = actordetails
	return actordetails, nil
}

func getData(url string, target interface{}) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to get data from url")
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		log.Println("not enough arguments")
		return
	}
	count := 0
	actor1 := os.Args[1]
	actor2 := os.Args[2]
	if actor1 == "" || actor2 == "" {
		log.Println("Source actor name and target actor name is required.")
		return
	}

	path, role, err := findDegreesofActor(actor1, actor2)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	fmt.Println("\n Degrees of Separation: ", (len(path)-1)/2)
	for i := 0; i < len(path)-2; i += 2 {
		fmt.Println((i/2)+1, "Movie: ", path[i+1])
		fmt.Println(role[count], " : ", path[i])
		fmt.Println(role[count+1], " : ", path[i+2])
		count++
	}
}
