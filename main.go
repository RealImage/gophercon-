package main

import (
	"log"
	"os"
)

func main() {

	// Usage:
	/*
		Compile the program using "go build main.go -o main.exe"
		Following is an example usage:

		[For Windows]
		./main.exe amitabh-bachchan robert-de-niro

		[For Others]
		Follow OS specific extensions in place of ".exe"
	*/

	// Command Line Arguments
	artistA := os.Args[1]
	artistB := os.Args[2]

	personA, err := FetchEntityDetails[Person](artistA)
	if err != nil {
		log.Println(err)
	}

	Separation(*personA, artistB)
}
