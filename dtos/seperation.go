package dtos

import (
	"github.com/sinhamanav030/challange2015/models"
)

type SeperationDegreeResponse struct {
	SeperationDegree int
	Path             []models.RelationNode
}
