package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	t := time.Now()
	defer func() {
		fmt.Println("Total Time taken: ", time.Since(t).Seconds())
	}()

	// Readin the Command Line Arguments
	artistA := os.Args[1]
	artistB := os.Args[2]

	personA, err := FetchEntityDetails[Person](artistA)
	if err != nil {
		log.Println(err)
	}

	// Search if related

	Separation(*personA, artistB)
}
