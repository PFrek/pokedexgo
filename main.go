package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/PFrek/pokedexgo/internal/pokeapi"
	"github.com/PFrek/pokedexgo/internal/pokecache"
)

type commandConfig struct {
	Next     *string
	Previous *string
	Cache    *pokecache.Cache
	Pokedex  map[string]pokeapi.PokemonResult
}

type command struct {
	name        string
	description string
	callback    func(*commandConfig, string) error
}

func printPrompt() {
	fmt.Print("pokedex > ")
}

func getInput(scanner *bufio.Scanner) string {
	if scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading scanner:", err)
	}
	return ""
}

func runCommand(input string, config *commandConfig) error {
	commandPart, arg, _ := strings.Cut(input, " ")

	validCommands := getValidCommands()
	command, ok := validCommands[commandPart]
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %s", input))
	}

	return command.callback(config, arg)
}

func getValidCommands() map[string]command {
	return map[string]command{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the names of the next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Find the pokemon in the specified location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch the specified pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "View the information of caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func commandPokedex(config *commandConfig, _ string) error {
	fmt.Println("Your Pokedex:")
	if len(config.Pokedex) == 0 {
		fmt.Println("[No entries found]")
		return nil
	}

	for name := range config.Pokedex {
		fmt.Printf("- %s\n", name)
	}
	return nil
}

func commandInspect(config *commandConfig, pokemonName string) error {
	if len(pokemonName) == 0 {
		return errors.New("pokemonName cannot be empty")
	}

	data, ok := config.Pokedex[pokemonName]
	if !ok {
		return errors.New(fmt.Sprintf("%s has not been caught yet", pokemonName))
	}

	printPokemonData(data)
	return nil
}

func printPokemonData(data pokeapi.PokemonResult) {
	fmt.Printf("Name: %s\n", data.Name)
	fmt.Printf("Height: %v\n", data.Height)
	fmt.Printf("Weight: %v\n", data.Weight)
	fmt.Println("Stats:")
	for _, stat := range data.Stats {
		fmt.Printf("- %s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range data.Types {
		fmt.Printf("- %s\n", t.Type.Name)
	}
}

func commandCatch(config *commandConfig, pokemonName string) error {
	if len(pokemonName) == 0 {
		return errors.New("pokemonName cannot be empty")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	result, err := pokeapi.GetPokemon(pokemonName, config.Cache)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get pokemon: %v", err))
	}

	caught := caughtPokemon(result.BaseExperience)

	if !caught {
		fmt.Printf("%s escaped!\n", pokemonName)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemonName)
	config.Pokedex[pokemonName] = *result

	return nil
}

func caughtPokemon(baseExp int) bool {
	target := 50 - 5*(baseExp/50)

	if target <= 0 {
		target = 1
	}

	roll := rand.Intn(100)
	return roll < target
}

func commandExplore(config *commandConfig, locationName string) error {
	if len(locationName) == 0 {
		return errors.New("locationName cannot be empty")
	}

	fmt.Printf("Exploring %s...\n", locationName)
	result, err := pokeapi.GetLocationPokemon(locationName, config.Cache)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get location pokemon: %v", err))
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range result {
		fmt.Printf("- %s\n", pokemon)
	}

	return nil
}

func commandMap(config *commandConfig, _ string) error {
	if config.Next == nil {
		return errors.New("Cannot go forward, already in last page")
	}

	result, err := pokeapi.GetLocations(config.Next, config.Cache)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get locations: %v", err))
	}

	config.Next = result.Next
	config.Previous = result.Previous

	pokeapi.PrintLocationNames(result.Locations)

	return nil
}

func commandMapBack(config *commandConfig, _ string) error {
	if config.Previous == nil {
		return errors.New("Cannot go back, already in first page")
	}
	result, err := pokeapi.GetLocations(config.Previous, config.Cache)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get locations: %v", err))
	}

	config.Next = result.Next
	config.Previous = result.Previous

	pokeapi.PrintLocationNames(result.Locations)

	return nil
}

func commandHelp(_ *commandConfig, _ string) error {
	validCommands := getValidCommands()
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range validCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()
	return nil
}

func commandExit(_ *commandConfig, _ string) error {
	fmt.Println("Exiting the Pokedex...")
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	initialNext := "https://pokeapi.co/api/v2/location-area/"
	config := commandConfig{
		Next:     &initialNext,
		Previous: nil,
		Cache:    pokecache.NewCache(5 * time.Minute),
		Pokedex:  make(map[string]pokeapi.PokemonResult),
	}

	for {
		printPrompt()
		textInput := getInput(scanner)

		err := runCommand(textInput, &config)
		if err != nil {
			fmt.Println("Error:", err)
		}

		if textInput == "exit" {
			break
		}
	}
}
