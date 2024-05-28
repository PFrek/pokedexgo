package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PFrek/pokedexgo/internal/pokecache"
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

func GetLocations(pageUrl *string, cache *pokecache.Cache) (*LocationsResult, error) {
	url := "https://pokeapi.co/api/v2/location-area/"
	if pageUrl != nil {
		url = *pageUrl
	}

	cachedValue, ok := cache.Get(url)
	if ok {
		fmt.Println("Cache hit, using cached value")
		return parseLocationsJson(cachedValue)
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

	cache.Add(url, body)

	return parseLocationsJson(body)
}

func GetLocationPokemon(locationName string, cache *pokecache.Cache) ([]string, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	cachedValue, ok := cache.Get(url)
	if ok {
		fmt.Println("Cache hit, using cached value")
		result, err := parseLocationJson(cachedValue)
		if err != nil {
			return nil, err
		}

		return extractPokemonNames(result), nil
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

	cache.Add(url, body)

	result, err := parseLocationJson(body)
	if err != nil {
		return nil, err
	}

	return extractPokemonNames(result), nil
}

func parseLocationsJson(body []byte) (*LocationsResult, error) {
	var result LocationsResult
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json parsing error: %v", err))
	}

	return &result, nil
}

func parseLocationJson(body []byte) (*LocationResult, error) {
	var result LocationResult
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json parsing error: %v", err))
	}

	return &result, nil
}

func extractPokemonNames(results *LocationResult) []string {
	pokemon := []string{}
	encounters := results.PokemonEncounters

	for _, encounter := range encounters {
		pokemon = append(pokemon, encounter.Pokemon.Name)
	}

	return pokemon
}

type LocationResult struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}
