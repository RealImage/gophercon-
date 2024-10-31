package main

import (
	"fmt"
	"os"
	"sync"
)

var (
	// semaphores and concurrency safe Cache
	sm    = NewSyncManager()
	Cache = sync.Map{}

	queue = []QueueData{} // A Queue to track artist details for BFS traversal

	parent        = make(map[string]Path) // Parent Map to keep track of the Path
	currentPerson = ""
)

func Separation(artistA Person, artistB string) {
	go func() {
		printResult(artistA.URL)
		os.Exit(0)
	}()

	queue = append(queue, QueueData{
		URL:      artistA.URL,
		Distance: 0,
	})

	for len(queue) > 0 {
		// Pop a node from Queue
		current := queue[0]
		queue = queue[1:]

		// Update the Path on each iteration
		key := current.URL
		if _, ok := parent[key]; !ok {
			parent[key] = Path{
				ParentURL:  current.ParentURL,
				Movie:      current.Movie,
				Role:       current.Role,
				ParentRole: current.ParentRole,
			}
		}

		if current.URL == artistB {
			currentPerson = current.URL
			sm.dos <- current.Distance
		}

		person, ok := Cache.Load(current.URL)
		if !ok {
			/* person isn't present in the cache.
			Fetch personDetails and update the cache.*/
			personDetails, err := FetchEntityDetails[Person](current.URL)
			if err != nil {
				if err.Error()[:4] == "403" {
					/* Storing the details when encountered 403 error
					so that we do not make a call to the same url again */
					Cache.Store(current.URL, Person{
						Type: "Forbidden",
					})
				}
			}
			if personDetails != nil {
				Cache.Store(current.URL, *personDetails)
				person = *personDetails
			}
		}

		if person != nil && person.(Person).Type != "Forbidden" {
			// Iterate through the MovieList to find related Persons
			for _, movie := range person.(Person).Movies {
				sm.wg.Add(1)
				// Function to handle the movie Data. i.e. find the linked artists and push them on the queue
				go handleMovieData(movie, current, artistB)
			}
			sm.wg.Wait()
		}
	}
}

func handleMovieData(m Details, current QueueData, artistB string) {
	defer sm.wg.Done()

	movie, ok := Cache.Load(m.URL)
	if !ok {
		/* Movie isn't present in the cache.
		Fetch movieDetails and update the cache	*/
		movieDetails, err := FetchEntityDetails[Movie](m.URL)
		if err != nil {
			if err.Error()[:4] == "403" {
				/* Storing the details when encountered 403 error
				so that we do not make a call to the same url again */
				Cache.Store(m.URL, Movie{
					Type: "Forbidden",
				})
			}
		}
		if movieDetails != nil {
			Cache.Store(m.URL, *movieDetails)
			movie = *movieDetails
		}
	}

	// Check if movie is Valid
	if movie != nil && movie.(Movie).Type != "Forbidden" {
		// Get the total list of related artists and append them to the queue with added distance
		artists := append(movie.(Movie).Cast, movie.(Movie).Crew...)
		sm.mu.Lock()
		for _, a := range artists {
			// Push artists on the queue
			queue = append(queue, QueueData{
				// Artist Details pushed on the queue
				URL:   a.URL,
				Movie: m.URL,
				Role:  a.Role,

				// Parent details pushed on the queue
				ParentURL:  current.URL,
				ParentRole: m.Role,

				// Increment the distance
				Distance: current.Distance + 1,
			})

			if a.URL == artistB {

				// If found, update the path and signal degrees of separation
				key := a.URL
				parent[key] = Path{
					ParentURL:  current.URL,
					Movie:      m.URL,
					Role:       a.Role,
					ParentRole: current.Role,
				}
				currentPerson = a.URL
				sm.dos <- current.Distance + 1
			}
		}
		sm.mu.Unlock()
	}
}

func printResult(sourceArtistURL string) {
	degrees := <-sm.dos
	fmt.Println("Distance of Separation: ", degrees)

	parentPerson := parent[currentPerson]
	for {
		defer func(parentPerson Path, currentPerson string, count int) {
			// Fetch the name of parent and person
			parentName, personName, movieName := GetPathDetails(parentPerson.ParentURL, currentPerson, parentPerson.Movie)
			fmt.Printf("\n%d. Movie: %s", count, movieName)
			fmt.Printf("\n%s: %s", parentPerson.ParentRole, parentName)
			fmt.Printf("\n%s: %s\n", parentPerson.Role, personName)
		}(parentPerson, currentPerson, degrees)

		currentPerson = parentPerson.ParentURL
		if currentPerson == sourceArtistURL {
			break
		}
		parentPerson = parent[currentPerson]
		degrees--
	}
}
