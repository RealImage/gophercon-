package dtos

import "github.com/sinhamanav030/challange2015/models"

type MovieAPIResponse struct {
	Name string         `json:"name"`
	Cast []models.Actor `json:"cast"`
	Crew []models.Actor `json:"crew"`
}

type ActorAPIResponse struct {
	Name   string         `json:"name"`
	Movies []models.Movie `json:"movies"`
}
