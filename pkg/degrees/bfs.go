package degrees

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
)

func BfsWithPath(start, target string) (*Path, error) {
	visited := &sync.Map{}
	queue := make([]Path, 0)
	queue = append(queue, Path{Nodes: []Node{}})
	visited.Store(start, true)

	for len(queue) > 0 {
		currentNodes := queue[0]
		queue = queue[1:]
		currentPerson := start
		if len(currentNodes.Nodes) != 0 {
			currentPerson = currentNodes.Nodes[len(currentNodes.Nodes)-1].Url2
		}

		body, err := getData(currentPerson)
		if err != nil {
			log.Printf("error getting data for %s. error: %s", currentPerson, err.Error())
			continue
		}

		var actor RespActorDto
		err = json.Unmarshal(body, &actor)
		if err != nil {
			log.Printf("error in unmarshalling actorDTO for %s. error: %s", currentPerson, err.Error())
			continue
		}

		var wg sync.WaitGroup
		movieChan := make(chan RespMovieDto, len(actor.Movies))

		for _, movie := range actor.Movies {
			wg.Add(1)
			go func(movie commonConfig) {
				defer wg.Done()

				if _, loaded := visited.LoadOrStore(movie.Url, true); loaded {
					return
				}

				body, err := getData(movie.Url)
				if err != nil {
					log.Printf("error getting data for %s. error: %s", movie.Url, err.Error())
					return
				}

				var movieDTO RespMovieDto
				err = json.Unmarshal(body, &movieDTO)
				if err != nil {
					log.Printf("error in unmarshalling movieDTO for %s. error: %s", movie.Url, err.Error())
					return
				}

				movieChan <- movieDTO
			}(movie)
		}

		go func() {
			wg.Wait()
			close(movieChan)
		}()

		var role1, person1 string
		for movieDTO := range movieChan {
			for _, p := range append(movieDTO.Cast, movieDTO.Crew...) {
				if p.Url == currentPerson {
					role1 = p.Role
					person1 = p.Name
					break
				}
			}

			for _, p := range append(movieDTO.Cast, movieDTO.Crew...) {
				if _, loaded := visited.LoadOrStore(p.Url, true); loaded {
					continue
				}

				newNodes := Path{
					Nodes: append(currentNodes.Nodes, Node{
						Movie:   movieDTO.Name,
						Role1:   role1,
						Person1: person1,
						Role2:   p.Role,
						Person2: p.Name,
						Url2:    p.Url,
					}),
				}

				if p.Url == target {
					return &newNodes, nil
				}
				queue = append(queue, newNodes)
			}
		}
	}

	return &Path{}, errors.New("unable to get a path")
}
