package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/sinhamanav030/challange2015/models"
	moviebuff "github.com/sinhamanav030/challange2015/utils/movieBuff"
	mutexMap "github.com/sinhamanav030/challange2015/utils/mutexMap"
)

type SeperationDegree interface {
	Find(source, dest models.Actor)
}

func NewSeperationDegreeService() SeperationDegree {
	return &seperationDegreeService{
		client: moviebuff.NewClient(),
		logger: *log.Default(),
	}
}

type seperationDegreeService struct {
	client moviebuff.Client
	logger log.Logger
}

func (s *seperationDegreeService) Find(source, dest models.Actor) {
	queue := []models.QueueNode{}
	visited := mutexMap.NewCallVisitedMutex()
	sourceNode := models.QueueNode{
		Person: source,
	}
	visited.Set(source.URL)
	queue = append(queue, sourceNode)

	degree := 1

	for len(queue) > 0 {
		ctx, cancelCtx := context.WithCancel(context.Background())

		queueWg := &sync.WaitGroup{}
		queueChan := make(chan models.QueueReader, 1)

		nodeJobChan := make(chan models.QueueNode, len(queue))
		nodeResChan := make(chan models.QueueNode)

		queueWg.Add(1)
		go func() {
			defer queueWg.Done()
			nxtLvlQueue := []models.QueueNode{}
			for node := range nodeResChan {
				if node.IsDest {
					queueChan <- models.QueueReader{DestFound: true, ResultNode: node}
					cancelCtx()
					return
				}
				nxtLvlQueue = append(nxtLvlQueue, node)
			}
			queueChan <- models.QueueReader{Queue: nxtLvlQueue}
		}()

		queueWg.Add(1)
		go func() {
			defer queueWg.Done()
			wg := &sync.WaitGroup{}
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					s.traverse(ctx, nodeResChan, nodeJobChan, source, dest, visited)
				}()
			}
			for _, node := range queue {
				nodeJobChan <- node
			}
			close(nodeJobChan)

			wg.Wait()

			close(nodeResChan)
		}()

		queueWg.Wait()
		close(queueChan)

		res := <-queueChan

		if res.DestFound {
			s.printRelation(res.ResultNode, degree)
			return
		}

		queue = res.Queue
		degree += 1
	}

	s.printRelation(models.QueueNode{}, 0)
}

func (s *seperationDegreeService) traverse(ctx context.Context, nodeResChn chan<- models.QueueNode, nodeJobChn <-chan models.QueueNode, source models.Actor, dest models.Actor, visited *mutexMap.CallVisitedMutex) {
	for curNode := range nodeJobChn {
		if len(curNode.Person.URL) == 0 {
			s.logger.Print("Skipping Invalid Node")
			continue
		}

		resp, err := s.client.FetchActor(curNode.Person.URL)
		if err != nil {
			s.logger.Print("error while fetching actor details")
			fmt.Println("here :", curNode.Person.URL)
			continue
		}

		for _, movie := range resp.Movies {
			if visited.Get(movie.URL) {
				continue
			}
			visited.Set(movie.URL)

			movieDetails, err := s.client.FetchMovie(movie.URL)
			if err != nil || movieDetails == nil {
				s.logger.Print("error while fetching movie details")
				continue
			}

			movieDetails.Cast = append(movieDetails.Cast, movieDetails.Crew...)
			for _, actor := range movieDetails.Cast {
				if strings.Compare(source.URL, curNode.Person.URL) == 0 && strings.Compare(curNode.Person.URL, actor.URL) == 0 {
					curNode.Person.Name = actor.Name
					curNode.Person.Role = actor.Role
				}

				if visited.Get(actor.URL) {
					continue
				}

				relationNode := models.RelationNode{
					FrstPerson:   curNode.Person,
					SecondPerson: actor,
					Movie:        movie.Name,
				}

				node := models.QueueNode{
					Person: actor,
				}

				if strings.Compare(dest.URL, actor.URL) == 0 {
					node.IsDest = true
				}

				node.Path = append(node.Path, curNode.Path...)
				node.Path = append(node.Path, relationNode)

				visited.Set(actor.URL)

				select {
				case nodeResChn <- node:
				case <-ctx.Done():
					return
				}

				if node.IsDest {
					return
				}

			}
		}
	}
}

func (s *seperationDegreeService) printRelation(node models.QueueNode, degree int) {
	if degree == 0 {
		fmt.Println("\nNo Relation found.")
		return
	}

	fmt.Println("\nDegree of Seperation: ", degree)
	for i, relation := range node.Path {
		fmt.Printf("%d. Movie: %s\n", i+1, relation.Movie)
		fmt.Printf("%s: %s\n", relation.FrstPerson.Role, relation.FrstPerson.Name)
		fmt.Printf("%s: %s\n", relation.SecondPerson.Role, relation.SecondPerson.Name)
		fmt.Println()
	}
}
