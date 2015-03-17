package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
	//"github.com/trustmaster/goflow"
)

type Degree struct {
	source       string
	target       string
	degree       int
	res          result
	graphIn      chan result
	out          chan bool
	nextUrlsChan chan []urls
	urlsChan     chan []urls
	personUrls   map[string]bool
	movieUrls    map[string]bool

	wg          sync.WaitGroup
	nextUrls    []urls
	connections []connection
}

func (d *Degree) FindDegree(src string, target string) {
	var err error
	for j := 0; j < 20; j++ {
		// Get the person's data
		_, err = getPersonData(src)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("Invalid person url")
		return
	}
	for j := 0; j < 20; j++ {
		// Get the person's data
		_, err = getPersonData(target)
		if err == nil {
			break
		}
	}
	if err != nil {
		fmt.Println("Invalid person url")
		return
	}
	d.graphIn = make(chan result)
	d.out = make(chan bool)
	d.nextUrlsChan = make(chan []urls)
	d.urlsChan = make(chan []urls)
	d.source = src
	d.target = target
	d.personUrls = make(map[string]bool)
	d.movieUrls = make(map[string]bool)
	r := result{}
	r.currentUrls = append(r.currentUrls, urls{src, nil})
	go d.handleInput()
	d.graphIn <- r
	ok := <-d.out
	fmt.Println("Success =", ok)
	fmt.Println("Degree of seperation:", d.degree)
	fmt.Println("Connections :")
	for i := range d.connections {
		fmt.Println(i+1, ".Movie:", d.connections[i].movie)
		fmt.Println(d.connections[i].firstRole, ":", d.connections[i].first)
		fmt.Println(d.connections[i].secondRole, ":", d.connections[i].second)
	}
}

func (d *Degree) handleMovie(url string, r *result, i int, pName string, pRole string) {
	//	fmt.Println(url)
	var m movie
	var err error
	var nextUrls []urls
	for j := 0; j < 20; j++ {
		// Get the movie data
		m, err = getMovieData(url)
		if err == nil {
			break
		}
	}
	if err != nil {

	} else {
		// Iterate through the list of casts
		for cst := range m.Cast {
			if d.personUrls[m.Cast[cst].Url] != true {
				// Sanity check(If the currentdegree is greater than degree no need to continue)
				if r.currentDegree < d.degree || d.degree == 0 {
					// Leave the current person, might lead to infinite loop if not checked
					if m.Cast[cst].Url != r.currentUrls[i].url {
						// If the current cast is the target
						if m.Cast[cst].Url == d.target {
							d.degree = r.currentDegree
							connections := r.currentUrls[i].connections
							connections = append(connections, connection{m.Name, pName, pRole, m.Cast[cst].Name, m.Cast[cst].Role})
							d.connections = connections
							// close the channel
							d.out <- true
							close(d.graphIn)
							break
						} else {
							connections := r.currentUrls[i].connections
							connections = append(connections, connection{m.Name, pName, pRole, m.Cast[cst].Name, m.Cast[cst].Role})
							nextUrls = append(nextUrls, urls{m.Cast[cst].Url, connections})
						}
					}
				}
			}
		}
		// Iterate through the list of crews
		for crw := range m.Crew {
			if d.personUrls[m.Crew[crw].Url] != true {
				// Sanity check(If the currentdegree is greater than degree no need to continue)
				if r.currentDegree < d.degree || d.degree == 0 {
					// Leave the current person, might lead to infinite loop if not checked
					if m.Crew[crw].Url != r.currentUrls[i].url {
						// If the current crew is the target
						if m.Crew[crw].Url == d.target {
							d.degree = r.currentDegree
							connections := r.currentUrls[i].connections
							connections = append(connections, connection{m.Name, pName, pRole, m.Crew[crw].Name, m.Crew[crw].Role})
							d.connections = connections
							// close the channel
							d.out <- true
							close(d.graphIn)
							break
						} else {
							connections := r.currentUrls[i].connections
							connections = append(connections, connection{m.Name, pName, pRole, m.Crew[crw].Name, m.Crew[crw].Role})
							nextUrls = append(nextUrls, urls{m.Crew[crw].Url, connections})
						}
					}
				}
			}
		}
	}
	//	fmt.Println(nextUrls)
	d.nextUrlsChan <- nextUrls
}

func (d *Degree) handlePerson(r *result, i int) {
	defer d.wg.Done()
	//	fmt.Println(r.currentUrls[i].url)
	var err error
	var per person
	var nextUrls []urls
	for j := 0; j < 20; j++ {
		// Get the person's data
		per, err = getPersonData(r.currentUrls[i].url)
		if err == nil {
			break
		}
	}
	if err != nil {

	} else {
		// Iterate through all his movies
		mvCount := 0
		for mv := range per.Movies {
			// Check if the movie is already traversed
			if d.movieUrls[per.Movies[mv].Url] != true {
				d.movieUrls[per.Movies[mv].Url] = true
				mvCount++
				go d.handleMovie(per.Movies[mv].Url, r, i, per.Name, per.Movies[mv].Role)
				//				var m movie
				//				for j := 0; j < 20; j++ {
				//					// Get the movie data
				//					m, err = getMovieData(per.Movies[mv].Url)
				//					if err == nil {
				//						break
				//					}
				//				}
				//				if err != nil {

				//				} else {
				//					// Iterate through the list of casts
				//					for cst := range m.Cast {
				//						// Sanity check(If the currentdegree is greater than degree no need to continue)
				//						if r.currentDegree < d.degree || d.degree == 0 {
				//							// Leave the current person, might lead to infinite loop if not checked
				//							if m.Cast[cst].Url != r.currentUrls[i].url {
				//								// If the current cast is the target
				//								if m.Cast[cst].Url == d.target {
				//									d.degree = r.currentDegree
				//									connections := r.currentUrls[i].connections
				//									connections = append(connections, connection{m.Name, per.Name, per.Movies[mv].Role, m.Cast[cst].Name, m.Cast[cst].Role})
				//									d.connections = connections
				//									// close the channel
				//									d.out <- true
				//									close(d.graphIn)
				//									break
				//								} else {
				//									connections := r.currentUrls[i].connections
				//									connections = append(connections, connection{m.Name, per.Name, per.Movies[mv].Role, m.Cast[cst].Name, m.Cast[cst].Role})
				//									d.nextUrls = append(d.nextUrls, urls{m.Cast[cst].Url, connections})
				//								}
				//							}
				//						}
				//					}
				//					// Iterate through the list of crews
				//					for crw := range m.Crew {
				//						// Sanity check(If the currentdegree is greater than degree no need to continue)
				//						if r.currentDegree < d.degree || d.degree == 0 {
				//							// Leave the current person, might lead to infinite loop if not checked
				//							if m.Crew[crw].Url != r.currentUrls[i].url {
				//								// If the current crew is the target
				//								if m.Crew[crw].Url == d.target {
				//									d.degree = r.currentDegree
				//									connections := r.currentUrls[i].connections
				//									connections = append(connections, connection{m.Name, per.Name, per.Movies[mv].Role, m.Crew[crw].Name, m.Crew[crw].Role})
				//									d.connections = connections
				//									// close the channel
				//									d.out <- true
				//									close(d.graphIn)
				//									break
				//								} else {
				//									connections := r.currentUrls[i].connections
				//									connections = append(connections, connection{m.Name, per.Name, per.Movies[mv].Role, m.Crew[crw].Name, m.Crew[crw].Role})
				//									d.nextUrls = append(d.nextUrls, urls{m.Crew[crw].Url, connections})
				//								}
				//							}
				//						}
				//					}
				//				}
			}
		}
		//		for j := 0; j < mvCount; j++ {
		//			tempUrls := <-d.nextUrlsChan
		//			nextUrls = append(nextUrls, tempUrls...)
		//		}
		//		fmt.Println(mvCount)
		for {
			if mvCount == 0 {
				break
			}
			select {
			case tempUrls, ok := <-d.nextUrlsChan:
				if !ok {
					break
				} else {
					mvCount--
					nextUrls = append(nextUrls, tempUrls...)
				}
			default:
				time.Sleep(1 * time.Microsecond)
			}
		}

	}
	d.urlsChan <- nextUrls
}

func (d *Degree) findDegree(r result) {
	r.currentDegree++
	noPerson := 0
	// Iterate through all the users
	for i := range r.currentUrls {
		// Check if the person is already traversed
		if d.personUrls[r.currentUrls[i].url] != true {
			noPerson++
			d.personUrls[r.currentUrls[i].url] = true
			d.wg.Add(1)
			go d.handlePerson(&r, i)
		}
	}

	d.nextUrls = nil
	for {
		if noPerson == 0 {
			break
		}
		select {
		case tempUrls, ok := <-d.urlsChan:
			if !ok {
				break
			} else {
				noPerson--

				d.nextUrls = append(d.nextUrls, tempUrls...)
			}
		default:
			time.Sleep(1 * time.Microsecond)
		}
	}
	//	for j := 0; j < noPerson; j++ {
	//		tempUrls := <-d.nextUrlsChan
	//		d.nextUrls = append(d.nextUrls, tempUrls...)
	//	}
	d.wg.Wait()

	if len(d.nextUrls) > 0 {
		// Send the url list again(Recursion)
		res := result{d.nextUrls, r.currentDegree}
		d.nextUrls = nil
		//		fmt.Println(res)
		d.graphIn <- res
	} else {
		d.out <- false
	}
}

// Recieves url list and processes it
func (d *Degree) handleInput() {
	for {
		// Get the url list
		select {
		case r, ok := <-d.graphIn:

			// Check if the channel is closed
			if !ok {
				d.out <- false
				break
			} else {
				go d.findDegree(r)
			}
		default:
			time.Sleep(1 * time.Microsecond)
		}
	}
}

func main() {
	// set the amount of CPU to be used
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS((numCPU * 3) / 4)
	// sanity check
	if len(os.Args) >= 3 {
		if os.Args[1] != os.Args[2] {
			d := new(Degree)
			start := time.Now()
			d.FindDegree(os.Args[1], os.Args[2])
			fmt.Println("Time taken to find degree =", time.Since(start).Seconds())
		}
	} else {
		fmt.Println("Please provide sufficient args")
	}
}
