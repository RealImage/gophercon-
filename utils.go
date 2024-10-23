package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func FetchData(moviebuffURL string) ([]byte, error) {
	res, err := http.Get(fmt.Sprintf("%s/%s", baseURL, moviebuffURL))
	if err != nil {
		log.Println("Error fetching data: ", err)
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading data: ", err)
		return nil, err
	}

	switch res.StatusCode {
	case 200:
		return data, nil

	default:
		log.Println("Invalid status code: ", res.StatusCode)
		return nil, fmt.Errorf("movie not found")
	}
}
