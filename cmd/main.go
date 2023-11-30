package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

type Word struct {
	word     string
	selected bool
	category string
}

type Board struct {
	words          []Word
	selectedAmount int
	selectedWords  []Word
	matchedWords   map[string][]Word
	attempts       []bool
}

func (b *Board) addWord(word string, category string) {
	b.words = append(b.words, Word{word, false, category})
}

func (b *Board) clickWord(word string) {
	for i, w := range b.words {
		if w.word == word {
			b.words[i].selected = !b.words[i].selected
			if b.words[i].selected {
				b.selectedAmount++
				b.selectedWords = append(b.selectedWords, w)
			} else {
				b.selectedAmount--
				if index := b.findSelectedWordIndex(w.word); index != -1 {
					b.deSelectWord(index)
				}
			}
		}
	}

	if b.selectedAmount == 4 {
		b.tryMatch()
		b.displayAttemptsAmount()
	}
}

func (b *Board) tryMatch() {
	category := b.selectedWords[0].category
	println(category)
	matchedWords := []Word{}
	for _, w := range b.selectedWords {
		if index := b.findWordIndex(w.word); index != -1 {
			if b.words[index].category != category {
				b.addAttempt(false)
				return
			}
			matchedWords = append(matchedWords, b.words[index])
		}
	}
	b.matchedWords[category] = matchedWords
	for _, w := range matchedWords {
		if index := b.findWordIndex(w.word); index != -1 {
			b.removeWord(index)
		}
	}
	b.addAttempt(true)
}

func (b *Board) addAttempt(success bool) {
	b.resetSelection()
	b.attempts = append(b.attempts, success)
}

func (b *Board) resetSelection() {
	b.selectedWords = []Word{}
	b.selectedAmount = 0
}

func (b *Board) deSelectWord(index int) {
	b.selectedWords = append(b.selectedWords[:index], b.selectedWords[index+1:]...)
}

func (b *Board) findWordIndex(word string) int {
	for i, w := range b.words {
		if w.word == word {
			return i
		}
	}
	return -1
}

func (b *Board) findSelectedWordIndex(word string) int {
	for i, w := range b.selectedWords {
		if w.word == word {
			return i
		}
	}
	return -1
}

func (b *Board) removeWord(index int) {
	b.words = append(b.words[:index], b.words[index+1:]...)
}

func (b *Board) addLine(line []string, category string) error {
	if len(line) != 4 {
		return errors.New("a line must have 4 words")
	}
	for _, word := range line {
		b.addWord(word, category)
	}
	return nil
}

func (b *Board) shuffleWords() {
	rand.Shuffle(len(b.words), func(i, j int) { b.words[i], b.words[j] = b.words[j], b.words[i] })
}

func (b *Board) displayBoard() {
	// Clear the screen
	fmt.Print("\033[H\033[2J")
	// Display matched words
	fmt.Println("\nMatched words:")
	b.displayMatchedWords()

	// Display the 4x4 grid
	fmt.Println("\nBoard:")
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			index := 4*i + j
			if index < len(b.words) {
				fmt.Printf("%2d: %-5s\t\t", index, b.words[index].word)
			}
		}
		fmt.Println()
	}
}

func (b *Board) displayMatchedWords() {
	for category, words := range b.matchedWords {
		fmt.Printf("%s: ", strings.ToUpper(category))
		for _, word := range words {
			fmt.Printf("[%s] ", strings.Title(strings.ToLower(word.word)))
		}
		fmt.Println()
	}
}

func (b *Board) displayAttemptsAmount() {
	fmt.Printf("\nAttempts: %d\n", len(b.attempts))
}

func (b *Board) displayAttemptsFull() {
	fmt.Printf("\nAttempts: %d\n\n", len(b.attempts))
	for _, v := range b.attempts {
		if v {
			fmt.Print("[\u2713] ")
		} else {
			fmt.Print("[X] ")
		}
	}
}

func assembleBoard(words map[string][]string) Board {
	b := Board{
		[]Word{},
		0,
		[]Word{},
		map[string][]Word{},
		[]bool{},
	}
	for category, line := range words {
		b.addLine(line, category)
	}
	return b
}

func gameLoop(b *Board) {
	scanner := bufio.NewScanner(os.Stdin)
	for len(b.words) > 0 {
		// Display the board
		b.displayBoard()

		// Get user input
		fmt.Print("\nEnter the index of a word to select: ")
		scanner.Scan()
		input := scanner.Text()

		// Convert the input to an integer
		index, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		// Check if the index is valid
		if index < 0 || index >= len(b.words) {
			fmt.Println("Invalid index. Please enter a number between 0 and", len(b.words)-1)
			continue
		}

		// Select the word
		b.clickWord(b.words[index].word)
	}

	// Declare the player a winner
	fmt.Println("\nCONGRATULATIONS, YOU'VE WON!")
	fmt.Println()
	b.displayMatchedWords()
	b.displayAttemptsFull()
}

func main() {
	words := map[string][]string{
		"animals": {"dog", "cat", "bird", "fish"},
		"colors":  {"red", "blue", "green", "yellow"},
		"fruits":  {"apple", "banana", "watermelon", "grape"},
		"shapes":  {"circle", "square", "triangle", "rectangle"},
	}
	board := assembleBoard(words)
	board.shuffleWords()
	gameLoop(&board)
}
