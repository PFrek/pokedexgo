package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		printPrompt()
		textInput := getInput(scanner)

		fmt.Println(textInput)
		if strings.ToLower(textInput) == "exit" {
			break
		}
	}
}
