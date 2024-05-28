package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationsResult struct {
	Count     int        `json:"count"`
	Next      *string    `json:"next"`
	Previous  *string    `json:"previous"`
	Locations []Location `json:"results"`
}

func PrintLocationNames(locations []Location) {
	for _, location := range locations {
		fmt.Println(location.Name)
	}
}

func GetLocations(pageUrl *string) (*LocationsResult, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if pageUrl != nil {
		url = *pageUrl
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Request error: %v", err))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Body parsing error: %v", err))
	}

	var result LocationsResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json parsing error: %v", err))
	}

	return &result, nil
}
