package services

import (
	"fmt"
	"log"
)

// FindDegreesOfSeparation performs BFS to find the shortest path
func FindDegreesOfSeparation(start, target string) (int, []string) {
	if start == target {
		return 0, []string{"identical person"}
	}

	bfsQueue := []string{start}
	visitedNode := make(map[string]bool)
	parentNode := make(map[string]string)

	visitedNode[start] = true

	for len(bfsQueue) > 0 {
		personID := bfsQueue[0]
		bfsQueue = bfsQueue[1:]

		// Fetch person details
		person, err := GetPersonData(personID)
		if err != nil {
			log.Printf("error fetching person details: %s\n", err)
			continue
		}

		// Process each movie the person has worked on
		for _, movie := range person.Movies {
			movie, err := GetMovieData(movie.URL) 
			if err != nil {
				log.Printf("error fetching movie details: %s\n", err)
				continue
			}

			// Process actors & directors (from cast & crew)
			for _, castMember := range movie.Cast {
				if visitedNode[castMember.URL] {
					continue
				}
				visitedNode[castMember.URL] = true
				parentNode[castMember.URL] = personID
				bfsQueue = append(bfsQueue, castMember.URL)

				if castMember.URL == target {
					return reconstructPath(parentNode, start, target)
				}
			}

			for _, crewMember := range movie.Crew {
				if visitedNode[crewMember.URL] {
					continue
				}
				visitedNode[crewMember.URL] = true
				parentNode[crewMember.URL] = personID
				bfsQueue = append(bfsQueue, crewMember.URL)

				if crewMember.URL == target {
					return reconstructPath(parentNode, start, target)
				}
			}
		}
	}

	return -1, []string{"no connection found"}
}

// reconstructPath builds the path from BFS parentNode map
func reconstructPath(parentNode map[string]string, start, target string) (int, []string) {
	path := []string{}
	current := target
	degree := 0

	for current != start {
		prev := parentNode[current]
		path = append([]string{fmt.Sprintf("%s → %s", prev, current)}, path...)
		current = prev
		degree++
	}

	return degree, path
}
