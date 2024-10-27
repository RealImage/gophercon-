package moviebuff

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sinhamanav030/challange2015/dtos"
)

type Client interface {
	FetchActor(actor string) (*dtos.ActorAPIResponse, error)
	FetchMovie(movie string) (*dtos.MovieAPIResponse, error)
}

type client struct {
	httpClient *http.Client
	logger     log.Logger
}

func NewClient() Client {
	client := &client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: *log.Default(),
	}
	return client
}

func (c *client) makeHttpReq(suffix string) (*http.Response, error) {
	httpUrl := fmt.Sprintf("https://data.moviebuff.com/%s", suffix)
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.httpClient.Do(req)
	return res, err
}

func (c *client) FetchActor(actor string) (*dtos.ActorAPIResponse, error) {
	res, err := c.makeHttpReq(actor)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data dtos.ActorAPIResponse
	if res.StatusCode != 200 {
		return &data, nil
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	return &data, err
}

func (c *client) FetchMovie(movie string) (*dtos.MovieAPIResponse, error) {
	res, err := c.makeHttpReq(movie)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data dtos.MovieAPIResponse
	if res.StatusCode != 200 {
		return &data, nil
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	return &data, err
}
