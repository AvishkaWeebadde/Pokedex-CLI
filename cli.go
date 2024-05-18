package main

import (
	"bufio"
	"fmt"
	"os"
)

// Initialize commandsMap in the init function
var commandsMap map[string]cliCommand

func init() {
	commandsMap = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Prints the list of available commands",
			callback: func() error {
				fmt.Println("List of available commands:")
				for name, command := range commandsMap {
					fmt.Printf("%s: %s\n", name, command.description)
				}
				return nil
			},
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex CLI",
			callback: func() error {
				fmt.Println("Goodbye!")
				os.Exit(0)
				return nil
			},
		},
	}

	fmt.Println("Welcome to the Pokedex CLI!")
	printDefaultPrompt()
}

func printDefaultPrompt() {
	fmt.Print("Pokedex CLI>")
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	for reader.Scan() {
		text := reader.Text()
		command, ok := commandsMap[text]
		if !ok {
			fmt.Println("Command not found. Type 'help' to see the list of available commands.")
			printDefaultPrompt()
			continue
		}
		err := command.callback()
		if err != nil {
			fmt.Printf("Error executing command: %s\n", err)
		}
		printDefaultPrompt()
	}
}
