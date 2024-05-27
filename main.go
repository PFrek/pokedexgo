package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	// "strings"
)

type command struct {
	name        string
	description string
	callback    func() error
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

func runCommand(input string) error {
	validCommands := getValidCommands()
	command, ok := validCommands[input]
	if !ok {
		return errors.New(fmt.Sprintf("invalid command %s", input))
	}

	return command.callback()
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
	}
}

func commandHelp() error {
	validCommands := getValidCommands()
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range validCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	fmt.Println()
	return nil
}

func commandExit() error {
	fmt.Println("Exiting the Pokedex...")
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		printPrompt()
		textInput := getInput(scanner)

		err := runCommand(textInput)
		if err != nil {
			fmt.Println("Error:", err)
		}

		if textInput == "exit" {
			break
		}
	}
}
