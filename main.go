package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)




type cast struct {
	Name string `json:"name"`
	Role string `json:"role"`
	URL  string `json:"url"`
}

type movie struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Role string `json:"role"`
}

type data struct {
	URL    string  `json:"url"`
	Type   string  `json:"type"`
	Name   string  `json:"name"`
	Movies []movie `json:"movies"`
	Casts  []cast  `json:"cast"`
	Crew   []cast  `json:"crew"`
}

type out struct {
	r   *relationship
	err error
}

type skippable struct {
	reason string
}

func (s skippable) Error() string {
	return s.reason
}

type relationship struct {
	cast1 cast
	cast2 cast
	movie string
	path  []relationship
}

type response struct {
	d   data
	err error
}

func validateInputs() bool {
	res := make(chan response)
	throttle := time.Tick(time.Second / 800)

	go getdata(os.Args[1], res, throttle)

	r := <-res
	if r.err != nil {
		return false
	}

	go getdata(os.Args[2], res, throttle)

	r = <-res
	if r.err != nil {
		return false
	}

	return true
}

func main() {
	if len(os.Args) < 3 {
		fmt.Print("Usage: degrees <actor-name> <actor-name>")
	}

	ok := validateInputs()
	if !ok {
		fmt.Print("Invalid ")
	}

	out := make(chan out)
	q := NewQueue()
	q.enqueue(workRequest{actor: os.Args[1]})

	go getrelationship(q, out)

	for {
		o := <-out
		if o.err != nil {
			if _, ok := o.err.(skippable); ok {
				continue
			}
			panic(o.err)
		}

		if o.r == nil {
			continue
		}

		if o.r.cast2.URL == os.Args[2] {
			var result []relationship
			result = o.r.path
			result = append(result, *o.r)

			fmt.Println("Degrees of seperation:", len(result))
			fmt.Println()

			for i, r := range result {
				fmt.Printf("%d. Movie: %s\n", i+1, r.movie)
				fmt.Println(r.cast1.Role+":", r.cast1.Name)
				fmt.Println(r.cast2.Role+":", r.cast2.Name)
				fmt.Println()
			}
			return
		}

		path := make([]relationship, len(o.r.path))
		copy(path, o.r.path)
		path = append(path, *o.r)

		q.enqueue(workRequest{actor: o.r.cast2.URL, path: path})
	}
}

func getdata(id string, res chan<- response, throttle <-chan time.Time) {
	fmt.Println("fetching", id)
	<-throttle
	resp, err := http.Get("http://data.moviebuff.com/" + id)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		res <- response{d: data{}, err: err}
		return
	}

	if resp.StatusCode == http.StatusForbidden {
		res <- response{d: data{}, err: skippable{reason: "forbidden"}}
		return
	}

	jsonBlob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res <- response{d: data{}, err: err}
		return
	}

	var data data
	err = json.Unmarshal(jsonBlob, &data)
	res <- response{d: data, err: err}
	return
}

func getrelationship(q *queue, output chan<- out) {
	visited := map[string]bool{}
	throttle := time.Tick(time.Second / 800)

	for {
		if q.empty() {
			continue
		}

		i := q.dequeue()
		if visited[i.actor] {
			output <- out{}
		}

		res := make(chan response)
		go getdata(i.actor, res, throttle)

		actorProfile := <-res
		if actorProfile.err != nil {
			if _, ok := actorProfile.err.(skippable); ok {
				output <- out{}
				continue
			}
			output <- out{err: actorProfile.err}
			continue
		}

		visited[i.actor] = true

		var count int
		urlVsMovie := map[string]movie{}
		for _, movie := range actorProfile.d.Movies {
			if visited[movie.URL] {
				continue
			}

			go getdata(movie.URL, res, throttle)

			urlVsMovie[movie.URL] = movie
			visited[movie.URL] = true
			count++
		}

		for k := 0; k < count; k++ {
			resp := <-res
			if resp.err != nil {
				if _, ok := resp.err.(skippable); ok {
					continue
				}
				output <- out{err: resp.err}
				return
			}

			for _, c := range resp.d.Casts {
				if c.URL == i.actor || visited[c.URL] {
					continue
				}

				visited[c.URL] = true

				output <- out{
					r: &relationship{
						cast1: cast{Name: actorProfile.d.Name, URL: actorProfile.d.URL, Role: urlVsMovie[resp.d.URL].Role},
						cast2: c,
						movie: urlVsMovie[resp.d.URL].Name,
						path:  i.path,
					},
				}
			}

			for _, c := range resp.d.Crew {
				if c.URL == i.actor || visited[c.URL] {
					continue
				}

				visited[c.URL] = true
				output <- out{
					r: &relationship{
						cast1: cast{Name: actorProfile.d.Name, URL: actorProfile.d.URL, Role: urlVsMovie[resp.d.URL].Role},
						cast2: c,
						movie: urlVsMovie[resp.d.URL].Name,
						path:  i.path,
					},
				}
			}
		}
	}
}