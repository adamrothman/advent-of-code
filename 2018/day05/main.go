package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func readInput(filename string) (string, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening config file %s for reading: %s", path, err)
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
		reacted := polymer
		for i, j := 0, 1; j < len(polymer); i, j = i+1, j+1 {
			one, two := int32(reacted[i]), int32(reacted[j])
			if abs32(one-two) == 32 {
				fmt.Printf("reaction between %c and %c\n", one, two)
			}
		}

		if reacted == polymer {
			return reacted
		}
		polymer = reacted
	}
}

func abs32(x int32) int32 {
	y := x >> 31
	return (x ^ y) - y
}

func main() {
	filename := "input.txt"

	polymer, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	result := react(polymer)
	fmt.Println("Reaction result:", result)
}
