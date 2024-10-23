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
		log.Println("Error fetching data.", "URL", moviebuffURL, "Error", err)
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Error reading data.", "URL", moviebuffURL, "Error", err)
		return nil, err
	}

	switch res.StatusCode {
	case 200:
		return data, nil

	default:
		log.Println("Invalid response.", "URL", moviebuffURL, "Status code", res.StatusCode, "Data", string(data))
		return nil, fmt.Errorf("invalid response. Status code: %d. Data: %s", res.StatusCode, string(data))
	}
}
