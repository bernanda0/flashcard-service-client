package main

import (
	"fmt"
)

const (
	baseURL   = "http://localhost:4444"
	tokenFile = "token.json"
	tokenDir  = ".flashcardapp"
)

func main() {
	for {
		fmt.Println("Welcome to the flash card app!")

		token, err := readTokenFromFile()
		if err != nil {
			fmt.Println("[You're not authenticated. Please signup/login]")
			fmt.Println("1. Signup")
			fmt.Println("2. Login")
			fmt.Println("3. Exit")

			var choice int
			fmt.Print("Enter your choice: ")
			_, err := fmt.Scanln(&choice)
			if err != nil {
				fmt.Println("Invalid input. Please enter a number.")
				continue
			}

			switch choice {
			case 1:
				if signup() {
					continue
				}
				continue
			case 2:
				if login() {
					continue
				}
				continue
			case 3:
				fmt.Println("Goodbye!")
				return
			default:
				fmt.Println("Invalid choice. Please select a valid option.")
			}
		}

		// Display menu and handle user actions
		displayMenu()
		var choice int
		fmt.Print("Enter your choice: ")
		_, err = fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			listDecks(token, token.UserID)
		case 2:
			addDeck(token.AccessToken, token.UserID)
		case 3:
			logout()
			fmt.Println("Goodbye!")
			return
		case 4:
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}
