package main

import (
	"fmt"
)

type Connection struct {
	Movie       string `json:"movie"`
	Person1     string `json:"person1"`
	Person1Role string `json:"person1_role"`
	Person2     string `json:"person2"`
	Person2Role string `json:"person2_role"`
}

func PrintResult(result []Connection, degree int64) {
	fmt.Println("Degree of Seperation: ", degree)

	for i, connection := range result {
		fmt.Printf("%d. Movie: %s\n", i+1, connection.Movie)
		fmt.Printf("%s: %s\n", connection.Person1Role, connection.Person1)
		fmt.Printf("%s: %s\n", connection.Person2Role, connection.Person2)
		fmt.Println()
	}
}
