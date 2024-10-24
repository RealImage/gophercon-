package main

import (
	"context"
	"sync"
)

// GraphSource represents a node in the graph to search for connections.
type GraphSource struct {
	// PersonURL is the URL of the person in the Moviebuff API.
	PersonURL string
	// Degree is the number of connections between the source person and current person.
	Degree int64
	// Connections are the details of the connections between the source person and current person.
	Connections []Collaboration
}

// GraphTraversalWorker traverses the nodes of the level provided in the "jobs" channel and writes the nodes of the next level to the "results" channel.
func GraphTraversalWorker(ctx context.Context, ctxCancel context.CancelFunc, destPerson string, vs *VisitedSources, jobs <-chan GraphSource, results chan<- GraphSource) {
	for job := range jobs {
		// Fetch the person data.
		personData, err := FetchPerson(job.PersonURL)
		if err != nil {
			continue
		}

		// Iterate over the movies of the person.
		for _, movie := range personData.Movies {
			// Skip if the movie has already been visited.
			if vs.Has(movie.URL) {
				continue
			}

			vs.Add(movie.URL)
			// Fetch the movie data.
			movieData, err := FetchMovie(movie.URL)
			if err != nil {
				continue
			}

			// movieMembers is the list of cast and crew of the movie.
			movieMembers := append(movieData.Cast, movieData.Crew...)
			// Iterate over the members of the movie.
			for _, member := range movieMembers {
				// Skip if the member has already been visited.
				if vs.Has(member.URL) {
					continue
				}

				// Store the details of the collaboration between the source person and the member.
				collab := append(job.Connections, Collaboration{
					Movie:       movieData.Name,
					Person1:     personData.Name,
					Person1Role: movie.Role,
					Person2:     member.Name,
					Person2Role: member.Role,
				})

				gs := GraphSource{
					PersonURL:   member.URL,
					Degree:      job.Degree + 1,
					Connections: collab,
				}

				// Check if the destination person has been found.
				// If so, print the connections, cancel the context to stop the workers, and return.
				if member.URL == destPerson {
					PrintCollaborations(gs.Connections, gs.Degree)
					ctxCancel()
					return
				}

				// Add the member to the visited list and write the node to the results channel.
				vs.Add(member.URL)
				select {
				case results <- gs:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

// SeperationDegrees calculates the minimum degrees of seperation i.e. the number of connections between two people.
// This is classical graph theory problem where we need to find the shortest path between two people.
// Each node in the graph represents a person and each edge represents a connection between two people which is a movie.
// However, the tricky part is we only know source and destination people, the rest of the people are hidden. The rest of
// the people will be explored at runtime.
// The BFS algorithm the most efficient in this case as we need to find the shortest path and the rest of the nodes are hidden.
// So, as we do in BFS, we will traverse the graph level by level and the moment we find the destination, we will return the
// number of connections.
func SeperationDegrees(srcPerson string, destPerson string) {
	// vs will store both visited people and movies.
	vs := NewVisitedSources()
	queue := []GraphSource{}

	// Add the source person to the queue.
	vs.Add(srcPerson)
	queue = append(queue, GraphSource{
		PersonURL:   srcPerson,
		Degree:      0,
		Connections: []Collaboration{},
	})

	// Traverse the graph level by level.
	for len(queue) > 0 {
		// ctx will be used to stop processing further when the destination is found.
		ctx, ctxCancel := context.WithCancel(context.Background())

		// jobs is the queue of people to be explored.
		jobs := make(chan GraphSource, len(queue))
		// results is the next level of people to be explored.
		results := make(chan GraphSource)

		// rwGroup is used to wait for the reader (reading from "results" channel) and writer (writing to "results" channel) to finish.
		rwGroup := &sync.WaitGroup{}
		// rwGroupChn is used to send the next level of people to be explored from the reader.
		rwGroupChn := make(chan []GraphSource, 1)

		// As the we don't know how many people are in the next level, the size of the
		// "results" channel is set to 1. That means we need to haver reader always reading from the channel.
		// Hence, we need to add a goroutine to read from the channel.
		rwGroup.Add(1)
		go func() {
			defer rwGroup.Done()

			newQueue := []GraphSource{}
			for {
				select {
				case result, ok := <-results:
					if !ok {
						rwGroupChn <- newQueue
						return
					}
					newQueue = append(newQueue, result)

				case <-ctx.Done():
					return
				}
			}
		}()

		// Add a goroutine to process the people in the queue.
		rwGroup.Add(1)
		go func() {
			defer rwGroup.Done()

			// As the number of people to be explored could be huge, we need to process them in batches.
			wg := &sync.WaitGroup{}
			for i := 0; i < 100; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					GraphTraversalWorker(ctx, ctxCancel, destPerson, vs, jobs, results)
				}()
			}

			// Add the people in the queue to the jobs channel.
			for _, gs := range queue {
				jobs <- gs
			}
			close(jobs)

			// Wait for the workers to finish. Then close the results channel so that the reader can finish.
			wg.Wait()
			close(results)
		}()

		rwGroup.Wait()
		close(rwGroupChn)

		select {
		case <-ctx.Done():
			return
		case queue = <-rwGroupChn:
		}
	}

	PrintCollaborations(nil, 0)
}
