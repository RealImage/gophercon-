package main

import (
	"fmt"
	"log"
	"os"
	"qube_assignment/services"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <actor1> <actor2>\nExample: %s amitabh-bachchan robert-de-niro", os.Args[0], os.Args[0])
	}

	actor1 := os.Args[1]
	actor2 := os.Args[2]

	degrees, path := services.FindDegreesOfSeparation(actor1, actor2)
	fmt.Printf("\nDegrees of Separation: %d\n\n", degrees)
	for i, step := range path {
		fmt.Printf("%d. %s\n", i+1, step)
	}
}
