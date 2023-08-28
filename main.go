package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	baseURL   = "http://localhost:4444"
	tokenFile = "token.json"
	tokenDir  = ".flashcardapp"
)

func main() {
	defer waitForAnyKey()
	for {
		token, err := readTokenFromFile()
		if err != nil { // todo checking the token file, because it can be exploited
			choice := printLoginMenu()
			switch choice {
			case 1:
				login()
				continue
			case 2:
				signup()
				continue
			case 3:
				exit()
				return
			default:
				printError("Invalid choice. Please select a valid option.")
				waitForAnyKey()
				continue
			}
		}

		// Display menu and handle user actions
		choice := printMainMenu()
		switch choice {
		case 1:
			for {
				decks := listDecks(token, token.UserID)
				choice := displayDecks(&decks)
				if choice == 0 {
					break
				} else if choice != -1 {
					for {
						deck_id := decks[choice-1].DeckID
						cards := listCards(token, deck_id)
						choice := displayCards(&cards)
						if choice == 0 {
							addCard(token, deck_id)
							continue
						} else {
							break
						}
					}

				}
			}
		case 2:
			addDeck(token, token.UserID)
		case 3:
			logout()
			continue
		case 4:
			exit()
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}

func waitForAnyKey() {
	fmt.Println("\nPress any key to continue...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadByte() // Read a single byte (any key)
}
