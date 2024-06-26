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

func parseLocationsJson(body []byte) (*LocationsResult, error) {
	var result LocationsResult
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json parsing error: %v", err))
	}

	return &result, nil
}

func GetLocationPokemon(locationName string, cache *pokecache.Cache) ([]string, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + locationName

	cachedValue, ok := cache.Get(url)
	if ok {
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

func GetPokemon(pokemonName string, cache *pokecache.Cache) (*PokemonResult, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	cachedValue, ok := cache.Get(url)
	if ok {
		return parsePokemonJson(cachedValue)
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

	return parsePokemonJson(body)

}

func parsePokemonJson(body []byte) (*PokemonResult, error) {
	var result PokemonResult
	err := json.Unmarshal(body, &result)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Json parsing error: %v", err))
	}

	return &result, nil
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
				Chance          int       `json:"chance"`
				ConditionValues []*string `json:"condition_values"`
				MaxLevel        int       `json:"max_level"`
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

type PokemonResult struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string    `json:"name"`
	Order         int       `json:"order"`
	PastAbilities []*string `json:"past_abilities"`
	PastTypes     []*string `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string  `json:"back_default"`
		BackFemale       *string `json:"back_female"`
		BackShiny        string  `json:"back_shiny"`
		BackShinyFemale  *string `json:"back_shiny_female"`
		FrontDefault     string  `json:"front_default"`
		FrontFemale      *string `json:"front_female"`
		FrontShiny       string  `json:"front_shiny"`
		FrontShinyFemale *string `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string  `json:"front_default"`
				FrontFemale  *string `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string  `json:"front_default"`
				FrontFemale      *string `json:"front_female"`
				FrontShiny       string  `json:"front_shiny"`
				FrontShinyFemale *string `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string  `json:"back_default"`
				BackFemale       *string `json:"back_female"`
				BackShiny        string  `json:"back_shiny"`
				BackShinyFemale  *string `json:"back_shiny_female"`
				FrontDefault     string  `json:"front_default"`
				FrontFemale      *string `json:"front_female"`
				FrontShiny       string  `json:"front_shiny"`
				FrontShinyFemale *string `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string  `json:"back_default"`
					BackFemale       *string `json:"back_female"`
					BackShiny        string  `json:"back_shiny"`
					BackShinyFemale  *string `json:"back_shiny_female"`
					FrontDefault     string  `json:"front_default"`
					FrontFemale      *string `json:"front_female"`
					FrontShiny       string  `json:"front_shiny"`
					FrontShinyFemale *string `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string  `json:"back_default"`
					BackFemale       *string `json:"back_female"`
					BackShiny        string  `json:"back_shiny"`
					BackShinyFemale  *string `json:"back_shiny_female"`
					FrontDefault     string  `json:"front_default"`
					FrontFemale      *string `json:"front_female"`
					FrontShiny       string  `json:"front_shiny"`
					FrontShinyFemale *string `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string  `json:"back_default"`
					BackFemale       *string `json:"back_female"`
					BackShiny        string  `json:"back_shiny"`
					BackShinyFemale  *string `json:"back_shiny_female"`
					FrontDefault     string  `json:"front_default"`
					FrontFemale      *string `json:"front_female"`
					FrontShiny       string  `json:"front_shiny"`
					FrontShinyFemale *string `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string  `json:"back_default"`
						BackFemale       *string `json:"back_female"`
						BackShiny        string  `json:"back_shiny"`
						BackShinyFemale  *string `json:"back_shiny_female"`
						FrontDefault     string  `json:"front_default"`
						FrontFemale      *string `json:"front_female"`
						FrontShiny       string  `json:"front_shiny"`
						FrontShinyFemale *string `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string  `json:"back_default"`
					BackFemale       *string `json:"back_female"`
					BackShiny        string  `json:"back_shiny"`
					BackShinyFemale  *string `json:"back_shiny_female"`
					FrontDefault     string  `json:"front_default"`
					FrontFemale      *string `json:"front_female"`
					FrontShiny       string  `json:"front_shiny"`
					FrontShinyFemale *string `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string  `json:"front_default"`
					FrontFemale      *string `json:"front_female"`
					FrontShiny       string  `json:"front_shiny"`
					FrontShinyFemale *string `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string  `json:"front_default"`
					FrontFemale      *string `json:"front_female"`
					FrontShiny       string  `json:"front_shiny"`
					FrontShinyFemale *string `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string  `json:"front_default"`
					FrontFemale  *string `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string  `json:"front_default"`
					FrontFemale      *string `json:"front_female"`
					FrontShiny       string  `json:"front_shiny"`
					FrontShinyFemale *string `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string  `json:"front_default"`
					FrontFemale  *string `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}
