package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
)

var minDegress atomic.Int64

func Degrees(person1 string, person2 string, curDegrees int64) (int64, bool) {
	if curDegrees >= minDegress.Load() {
		return 0, false
	}

	person1Data, err := FetchPersonDetails(person1)
	if err != nil {
		return -1, false
	}

	visitedSources.Add(person1)
	movies := make([]Movie, 0)
	for _, movieRole := range person1Data.Movies {
		if visitedSources.Has(movieRole.URL) {
			continue
		}

		movieData, err := FetchMovieDetails(movieRole.URL)
		if err != nil {
			continue
		}

		visitedSources.Add(movieRole.URL)
		movies = append(movies, movieData)
	}

	for _, movie := range movies {
		for _, person := range movie.Cast {
			if person.URL == person2 {
				return curDegrees + 1, true
			}
		}

		for _, person := range movie.Crew {
			if person.URL == person2 {
				return curDegrees + 1, true
			}
		}
	}

	var wg sync.WaitGroup
	for _, movie := range movies {
		for _, person := range movie.Cast {
			if visitedSources.Has(person.URL) {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				degrees, ok := Degrees(person.URL, person2, curDegrees+1)
				if ok && degrees < minDegress.Load() {
					minDegress.Store(degrees)
				}
			}()
		}

		for _, person := range movie.Crew {
			if visitedSources.Has(person.URL) {
				continue
			}

			wg.Add(1)
			go func() {
				defer wg.Done()

				degrees, ok := Degrees(person.URL, person2, curDegrees+1)
				if ok && degrees < minDegress.Load() {
					minDegress.Store(degrees)
				}
			}()
		}
	}

	wg.Wait()
	return 0, false
}

func main() {
	minDegress.Store(math.MaxInt64)
	Degrees("amitabh-bachchan", "robert-de-niro", 0)
	fmt.Println(minDegress.Load())
}
