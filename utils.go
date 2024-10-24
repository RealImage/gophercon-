package main

import (
	"fmt"
)

// Collaboration represents a collaboration between two people in a movie.
type Collaboration struct {
	Movie       string
	Person1     string
	Person1Role string
	Person2     string
	Person2Role string
}

// PrintCollaborations prints a list of collaborations.
func PrintCollaborations(collabs []Collaboration, degree int64) {
	if len(collabs) == 0 {
		fmt.Println("No collaborations found.")
		return
	}

	fmt.Println("Degree of Seperation: ", degree)
	for i, connection := range collabs {
		fmt.Printf("%d. Movie: %s\n", i+1, connection.Movie)
		fmt.Printf("%s: %s\n", connection.Person1Role, connection.Person1)
		fmt.Printf("%s: %s\n", connection.Person2Role, connection.Person2)
		fmt.Println()
	}
}
