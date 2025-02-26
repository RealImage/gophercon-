package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Kindly provide correct arguments")
		return
	}
	FindDegreeofSeparation(args[1], args[2])
}

// To change people, update in MakeFile
// I am running the program once only, so not storing anything for the next call
