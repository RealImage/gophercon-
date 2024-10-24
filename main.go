package main

import (
	"fmt"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Invalid number of arguments.")
		fmt.Println("Usage: go run . <moviebuff_url_1> <moviebuff_url_2>")
		return
	}

	moviebuffURL1 := os.Args[1]
	moviebuffURL2 := os.Args[2]
	if moviebuffURL1 == moviebuffURL2 {
		fmt.Println("Error: Moviebuff URLs must be different.")
		return
	}

	fmt.Println("Searching Moviebuff...")
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	defer s.Stop()

	SeperationDegrees(moviebuffURL1, moviebuffURL2)
}
