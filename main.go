package main

import (
	"bufio"
	"errors"
	"fmt"
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
	}
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
