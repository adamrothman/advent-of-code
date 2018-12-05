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
		return nil, fmt.Errorf("opening config file %s for reading: %s", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	changes := make([]int64, 0)

	for scanner.Scan() {
		line := scanner.Text()
		i, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Printf("Error parsing int64 from %s: %s", line, err)
			continue
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
