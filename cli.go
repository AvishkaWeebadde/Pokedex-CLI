package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type LocationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var (
	commandsMap map[string]cliCommand
	offset      int
	limit       int = 20
)


/**
 * The init function is called before the main function and is used to initialize the CLI.
 * It initializes the offset variable to 0 and the limit variable to 20.
 * It also initializes the commandsMap variable, which is a map of strings to cliCommand structs.
 * The cliCommand struct has three fields: name, description, and callback.
 * The name field is the name of the command, the description field is a short description of the command, 
 * and the callback field is a function that is called when the command is executed.
 * The init function also prints a welcome message and the default prompt.
 */
func init() {

	offset = 0

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
		"map": {
			name:        "map",
			description: "Displays the names of next 20 location areas",
			callback: func() error {
				url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=%d&offset=%d", limit, offset)
				resp, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				var locationAreas LocationAreas
				err = json.NewDecoder(resp.Body).Decode(&locationAreas)
				if err != nil {
					log.Fatal(err)
				}

				for _, locationArea := range locationAreas.Results {
					fmt.Println(locationArea.Name)
				}

				offset += limit

				return nil
			},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of previous 20 location areas",
			callback: func() error {
				if offset > limit {
					offset -= 2 * limit
				} else {
					offset = 0
				}

				url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area?limit=%d&offset=%d", limit, offset)
				resp, err := http.Get(url)
				if err != nil {
					log.Fatal(err)
				}
				defer resp.Body.Close()

				var locationAreas LocationAreas
				err = json.NewDecoder(resp.Body).Decode(&locationAreas)
				if err != nil {
					log.Fatal(err)
				}

				for _, locationArea := range locationAreas.Results {
					fmt.Println(locationArea.Name)
				}

				offset += limit 

				return nil
			},
		},
	}

	fmt.Println("Welcome to the Pokedex CLI!")
	printDefaultPrompt()
}

/**
 * The printDefaultPrompt function prints the default prompt for the CLI, which is "Pokedex CLI>".
 * The prompt is printed in blue color using ANSI escape codes.
 */
func printDefaultPrompt() {
	const blue = "\033[34m"
	const reset = "\033[0m"
	fmt.Printf("%sPokedex CLI>%s", blue, reset)
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
