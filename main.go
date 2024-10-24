package main

import (
	"sync"
)

var visitedSources = NewVisitedSources()

func Degrees(person1 string, person2 string) {
	jobQu := []Job{
		{
			person:      person1,
			degree:      0,
			connections: []Collaboration{},
		},
	}
	visitedSources.Add(person1)

	for len(jobQu) > 0 {
		qusize := len(jobQu)
		jobs := make(chan Job, qusize)
		results := make(chan Job)
		poolGroup := &sync.WaitGroup{}
		poolResult := make(chan struct {
			jobQu []Job
			done  bool
		}, 1)

		poolGroup.Add(1)
		go func() {
			defer poolGroup.Done()

			workerGroup := &sync.WaitGroup{}
			for i := 0; i < 100; i++ {
				workerGroup.Add(1)
				go Worker(workerGroup, person2, jobs, results)
			}

			for i := 0; i < qusize; i++ {
				jobs <- jobQu[i]
			}
			close(jobs)

			workerGroup.Wait()
			close(results)
		}()

		poolGroup.Add(1)
		go func() {
			defer poolGroup.Done()

			newJobQ := []Job{}
			for result := range results {
				if result.done {
					PrintCollaborations(result.connections, result.degree)
					poolResult <- struct {
						jobQu []Job
						done  bool
					}{
						done: true,
					}
					return
				}

				newJobQ = append(newJobQ, result)
			}

			poolResult <- struct {
				jobQu []Job
				done  bool
			}{
				jobQu: newJobQ,
			}
		}()

		poolGroup.Wait()
		close(poolResult)

		res := <-poolResult
		if res.done {
			return
		}

		jobQu = res.jobQu
	}
}

func main() {
	Degrees("amitabh-bachchan", "robert-de-niro")
}
