package main

import (
	"log"
	"os"
	"strings"

	"github.com/sinhamanav030/challange2015/models"
	service "github.com/sinhamanav030/challange2015/service"
)

func main() {

	args := os.Args[1:]

	if len(args) < 2 || strings.Compare(args[0], args[1]) == 0 {
		log.Fatal("Incorrect: require unique source actor and dest actor url")
	}

	svc := service.NewSeperationDegreeService()

	svc.Find(models.Actor{URL: args[0]}, models.Actor{URL: args[1]})

}
