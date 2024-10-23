package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
	"os"
    // "time"
)

type Movie struct {
    Name string `json:"name"`
    URL  string `json:"url"`
    Role string `json:"role"`
    Cast []Person `json:"cast"`
    Crew []Person `json:"crew"`
}

type Person struct {
    URL    string  `json:"url"`
    Type   string  `json:"type"`
    Name   string  `json:"name"`
	Role   string  `json:"role"`
    Movies []Movie `json:"movies"`
}

var personCache = make(map[string]Person)
var movieCache = make(map[string]Movie)

func fetchJSON(url string, target interface{}) error {
    // time.Sleep(200 * time.Millisecond)

    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return errors.New("failed to fetch data")
    }

    return json.NewDecoder(resp.Body).Decode(target)
}

func getPerson(moviebuffURL string) (Person, error) {
    if person, exists := personCache[moviebuffURL]; exists {
        return person, nil
    }
    
    url := fmt.Sprintf("https://data.moviebuff.com/%s", moviebuffURL)
    var person Person
    if err := fetchJSON(url, &person); err != nil {
        return Person{}, err
    }
    
    personCache[moviebuffURL] = person
    return person, nil
}

func getMovie(moviebuffURL string) (Movie, error) {
    if movie, exists := movieCache[moviebuffURL]; exists {
        return movie, nil
    }
    
    url := fmt.Sprintf("https://data.moviebuff.com/%s", moviebuffURL)
    var movie Movie
    if err := fetchJSON(url, &movie); err != nil {
        return Movie{}, err
    }
    
    movieCache[moviebuffURL] = movie
    return movie, nil
}

type Node struct {
    Name    string
    Path    []string
}



func findDegrees(source, target string) ([]string,[]string,error) {
	role := []string{"Actor"}
    queue := []Node{{Name: source, Path: []string{source}}}
    visited := make(map[string]bool)
    visited[source] = true

    for len(queue) > 0 {
        currentNode := queue[0]
        queue = queue[1:]

        person, err := getPerson(currentNode.Name)
        if err != nil {
			// fmt.Println(currentNode.Name)
            continue
        }

        for _, movie := range person.Movies {
            movieData, err := getMovie(movie.URL)
            if err != nil {
                continue
            }


            for _, cast := range append(movieData.Cast, movieData.Crew...) {
                if cast.URL == target {
                    return append(currentNode.Path, movieData.Name, target),append(role, cast.Role), nil
                }
            }

            for _, actor := range append(movieData.Cast,movieData.Crew...)  {
                if !visited[actor.URL] {
                    visited[actor.URL] = true
                    queue = append(queue, Node{Name: actor.URL, Path: append(currentNode.Path, movieData.Name, actor.URL)})
                }
            }

        }
    }

    return nil,[]string{},errors.New("no connection found")
}


func main() {
    if len(os.Args) < 3 {
        fmt.Println("Usage: go run main.go <source> <target>")
        return
    }

    // Get source and target from command-line arguments
    source := os.Args[1]
    target := os.Args[2]
	if source == "" || target == ""{
		fmt.Println("Source actor name and target actor name is required.")
        return
	}
    
    // Find degrees of separation
    path,role, err := findDegrees(source, target)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
	fmt.Println(path)
	fmt.Println(role)
    // Output the result
	role_count:=0
    fmt.Println("\n Degrees of Separation: ", (len(path)-1)/2)
    for i := 0; i < len(path)-2; i += 2 {
		fmt.Println((i/2)+1,"Movie: ",path[i+1])
		fmt.Println(role[role_count]," : ",path[i])
		fmt.Println(role[role_count+1]," : ",path[i+2])
		role_count++

    }
}
