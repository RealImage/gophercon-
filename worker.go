package main

import "sync"

type Job struct {
	person      string
	degree      int64
	connections []Connection
	done        bool
}

func Worker(wg *sync.WaitGroup, person2 string, jobs <-chan Job, results chan<- Job) {
	defer wg.Done()

	for j := range jobs {
		personData, err := FetchPersonDetails(j.person)
		if err != nil {
			return
		}

		for _, movieRole := range personData.Movies {
			movieData, err := FetchMovieDetails(movieRole.URL)
			if err != nil {
				continue
			}

			people := append(movieData.Cast, movieData.Crew...)

			for _, person := range people {
				if visitedSources.Has(person.URL) {
					continue
				}

				connections := append(j.connections, Connection{
					Movie:       movieData.Name,
					Person1:     personData.Name,
					Person1Role: movieRole.Role,
					Person2:     person.Name,
					Person2Role: person.Role,
				})

				newJob := Job{
					person:      person.URL,
					degree:      j.degree + 1,
					connections: connections,
				}

				if person.URL == person2 {
					newJob.done = true
					results <- newJob
					return
				}

				visitedSources.Add(person.URL)
				results <- newJob
			}
		}
	}
}
