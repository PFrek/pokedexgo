package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PFrek/pokedexgo/internal/pokeapi"
	"os"
)

type commandConfig struct {
	Next     *string
	Previous *string
}

type command struct {
	name        string
	description string
	callback    func(*commandConfig) error
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
	validCommands := getValidCommands()
	command, ok := validCommands[input]
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %s", input))
	}

	return command.callback(config)
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
	}
}

func commandMap(config *commandConfig) error {
	if config.Next == nil {
		return errors.New("Cannot go forward, already in last page")
	}

	result, err := internal.GetLocations(config.Next)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get locations: %v", err))
	}

	config.Next = result.Next
	config.Previous = result.Previous

	internal.PrintLocationNames(result.Locations)

	return nil
}

func commandMapBack(config *commandConfig) error {
	if config.Previous == nil {
		return errors.New("Cannot go back, already in first page")
	}
	result, err := internal.GetLocations(config.Previous)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to get locations: %v", err))
	}

	config.Next = result.Next
	config.Previous = result.Previous

	internal.PrintLocationNames(result.Locations)

	return nil
}

func commandHelp(_ *commandConfig) error {
	validCommands := getValidCommands()
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range validCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()
	return nil
}

func commandExit(_ *commandConfig) error {
	fmt.Println("Exiting the Pokedex...")
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	initialNext := "https://pokeapi.co/api/v2/location-area/"
	config := commandConfig{
		Next:     &initialNext,
		Previous: nil,
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
