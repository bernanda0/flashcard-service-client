package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/eiannone/keyboard"
)

func header() {
	clearScreen()
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Println("| 📇                    Card Repetition App                      |")
	fmt.Println("+----------------------------------------------------------------+")
}

func msgNotAuthed() {
	fmt.Println("|         You're not authenticated. Please signup/login          |")
	fmt.Println("+----------------------------------------------------------------+")
}

func msgWelcome() {
	fmt.Println("|           Welcome! Lets study better using flashcard           |")
	fmt.Println("+----------------------------------------------------------------+")
}

func printLoginMenu() int {
	header()
	msgNotAuthed()
	fmt.Println("|    Welcome! Please select an option                            |")
	fmt.Println("|    1. Login                                                    |")
	fmt.Println("|    2. Signup                                                   |")
	fmt.Println("|    3. Exit                                                     |")
	fmt.Println("+----------------------------------------------------------------+")
	return chooseNumber(1, 3)
}

func printMainMenu() int {
	header()
	msgWelcome()
	fmt.Println("|    Main Menu:                                                  |")
	fmt.Println("|    1. My Decks                                                 |")
	fmt.Println("|    2. Add Decks                                                |")
	fmt.Println("|    3. Logout                                                   |")
	fmt.Println("|    4. Exit                                                     |")
	fmt.Println("+----------------------------------------------------------------+")
	return chooseNumber(1, 4)
}

func loginForm() (email, password string) {
	header()
	fmt.Println("|                   Fill your email & password                   |")
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Print("  > Email 	[xxxxx@xxx.xxx]	: ")
	fmt.Scanln(&email)
	fmt.Print("  > Password	[ 6 char ]	: ")
	fmt.Scanln(&password)

	return
}

func exit() {
	header()
	fmt.Println("|                   Thanks, see you again 🖐️                      |")
	fmt.Println("+----------------------------------------------------------------+")

}

func addCardForm() (question, answer string) {
	for {
		header()
		fmt.Println("|                 Fill the Q & A for your card                   |")
		fmt.Println("|                     [Fill 0 to go back]                        |")
		fmt.Println("+----------------------------------------------------------------+")
		fmt.Print("  > Question: ")

		reader := bufio.NewReader(os.Stdin)

		question, _ = reader.ReadString('\n')
		question = strings.TrimSpace(question)
		if question == string('0') {
			return "", ""
		}

		fmt.Print("  > Answer  : ")
		answer, _ = reader.ReadString('\n')
		answer = strings.TrimSpace(answer)
		if answer == string('0') {
			return "", ""
		}

		fmt.Println("+----------------------------------------------------------------+")
		fmt.Println("  > Question:", question)
		fmt.Println("  > Answer  :", answer)
		fmt.Println("+----------------------------------------------------------------+")
		fmt.Println("Confirm (Y/N):")

		char, _, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if char == 'Y' || char == 'y' {
			break
		}
	}
	return question, answer
}

func displayDecks(decks *[]DeckResponse) int {
	header()
	fmt.Println("|  What subject you want to recall. Choose deck to see the cards |")
	fmt.Println("|                       Choose 0 to back                         |")
	fmt.Println("+----------------------------------------------------------------+")
	fmt.Printf("| %-3s | %-39s | %-14s |\n", "No.", "Title", "Created At")
	fmt.Println("+----------------------------------------------------------------+")
	for i, deck := range *decks {
		createdAt := ""
		if deck.CreatedAt.Valid {
			createdAt = deck.CreatedAt.Time.Format("2006-01-02")
		}
		fmt.Printf("| %-3d | %-39s | %-14s |\n", i+1, deck.Title, createdAt)
	}
	fmt.Println("+----------------------------------------------------------------+")
	return chooseNumber(0, len(*decks))
}

func displayCards(cards *[]CardResponse) int {
	currentCardIndex := 0
	showAnswer := false
	archived := false
	card_size := len(*cards)
	for {
		header()
		fmt.Println("|          You can see the answer or do next/prev cards          |")
		fmt.Println("|                   + Press A to add card +                      |")
		fmt.Println("+----------------------------------------------------------------+")

		if card_size == 0 {
			fmt.Println("|               No card found! You can add first.                |")
		} else {
			card := (*cards)[currentCardIndex]

			fmt.Printf("| < PREV                   [ID : %-2d]                      NEXT > |\n", card.FlashcardID)
			fmt.Println("+----------------------------------------------------------------+")
			fmt.Println()
			fmt.Printf("  ❔  Question	: %s\n", card.Question)
			if showAnswer {
				fmt.Printf("  🙋  Answer	: %s\n", card.Answer)
			} else if archived {
				fmt.Println("<Archived>")
			}

		}
		fmt.Println()
		fmt.Println("+----------------------------------------------------------------+")
		fmt.Println("  ENTER		" + toggleShowAnswer(showAnswer))
		fmt.Println("  ESC		Back")

		_, key, err := keyboard.GetSingleKey()
		if err != nil {
			fmt.Println("Error:", err)
			return -1
		}

		switch key {
		case keyboard.KeyInsert:
			return 0
		case keyboard.KeyArrowRight:
			currentCardIndex = nextIndex(currentCardIndex, card_size)
		case keyboard.KeyArrowLeft:
			currentCardIndex = prevIndex(currentCardIndex, card_size)
		case keyboard.KeyEnter:
			showAnswer = !showAnswer
		case keyboard.KeyEsc:
			return -1
		}
	}
}

func toggleShowAnswer(showAnswer bool) string {
	if showAnswer {
		return "Hide Answer"
	}
	return "Show Answer"
}

func toggleArchive(archived bool) string {
	if archived {
		return "Unarchive"
	}
	return "Archive"
}

func nextIndex(index, length int) int {
	if length == 0 {
		return 0
	}
	return (index + 1) % length
}

func prevIndex(index, length int) int {
	if length == 0 {
		return 0
	}
	return (index - 1 + length) % length
}

func printError(err string) {
	fmt.Printf("\n ❗  %s\n", err)
}

func printSuccess(succ string) {
	fmt.Printf("\n ✔️  %s\n", succ)
}

func chooseNumber(a, b int) int {
	var choice int
	fmt.Print("  > Enter your choice: ")
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < a || choice > b {
		return -1
	}
	return choice
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
