package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func readInput(filename string) ([]int64, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	changes := make([]int64, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		i, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parsing int64: %s", err)
		}
		changes = append(changes, i)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return changes, nil
}

func calculateFrequency(changes []int64) (frequency int64) {
	for _, change := range changes {
		frequency += change
	}
	return
}

func calculateRepeatFrequency(changes []int64) (frequency int64) {
	frequencies := make(map[int64]bool)
	frequencies[frequency] = true

	for {
		for _, change := range changes {
			frequency += change
			if frequencies[frequency] {
				return
			}
			frequencies[frequency] = true
		}
	}
}

func main() {
	filename := "input.txt"

	changes, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	frequency := calculateFrequency(changes)
	fmt.Println("Single-pass frequency:", frequency)

	repeatFrequency := calculateRepeatFrequency(changes)
	fmt.Println("First repeated frequency:", repeatFrequency)
}
