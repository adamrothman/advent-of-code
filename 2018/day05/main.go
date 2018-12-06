package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func readInput(filename string) (string, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var text string

	for scanner.Scan() {
		text = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("reading input file: %s", err)
	}

	return text, nil
}

func react(polymer string) string {
	for {
		var result strings.Builder

		for i := 0; i < len(polymer); {
			current := polymer[i]

			if i+1 < len(polymer) {
				next := polymer[i+1]
				if areReactive(current, next) {
					i += 2
					continue
				}
			}

			result.WriteByte(current)
			i++
		}

		newPolymer := result.String()
		if newPolymer == polymer {
			return polymer
		}
		polymer = newPolymer
	}
}

func areReactive(x, y byte) bool {
	return abs16(int16(x)-int16(y)) == 32
}

func abs16(x int16) int16 {
	y := x >> 15
	return (x ^ y) - y
}

var letters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func findShortestAfterSingleExcision(polymer string) (shortest string, removedUnits string) {
	minLength := len(polymer)

	for _, letter := range letters {
		upper := strings.ToUpper(letter)

		replacer := strings.NewReplacer(letter, "", upper, "")
		excised := replacer.Replace(polymer)

		result := react(excised)
		resultLen := len(result)

		if resultLen < minLength {
			minLength = resultLen
			shortest = result
			removedUnits = letter + upper
		}
	}

	return
}

func main() {
	filename := "input.txt"

	polymer, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	result := react(polymer)
	fmt.Printf("Resulting polymer has %d units\n", len(result))

	shortest, removedUnits := findShortestAfterSingleExcision(polymer)
	fmt.Printf("Removing %s produced the shortest reacted polymer at %d units\n", removedUnits, len(shortest))
}
