package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func readInput(filename string) ([]string, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	boxIDs := make([]string, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		boxIDs = append(boxIDs, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return boxIDs, nil
}

func calculateChecksum(boxIDs []string) int64 {
	var doubles, triples int64

	for _, boxID := range boxIDs {
		counter := make(map[rune]int)

		for _, char := range boxID {
			counter[char]++
		}

		hasDouble, hasTriple := false, false

		for _, count := range counter {
			if hasDouble && hasTriple {
				break
			}

			if !hasDouble && count == 2 {
				hasDouble = true
			} else if !hasTriple && count == 3 {
				hasTriple = true
			}
		}

		if hasDouble {
			doubles++
		}
		if hasTriple {
			triples++
		}
	}

	return doubles * triples
}

func findSimilarBoxIDs(boxIDs []string) (string, string, int) {
	for firstIndex, firstID := range boxIDs {
		for secondIndex, secondID := range boxIDs {
			if firstIndex == secondIndex {
				continue
			}

			differences := make([]int, 0)

			for i := 0; i < len(firstID); i++ {
				if firstID[i] != secondID[i] {
					differences = append(differences, i)
				}
			}

			if len(differences) == 1 {
				return firstID, secondID, differences[0]
			}
		}
	}

	return "", "", 0
}

func main() {
	filename := "input.txt"

	boxIDs, err := readInput(filename)
	if err != nil {
		log.Fatalf("Error reading input from %s: %s\n", filename, err)
	}

	checksum := calculateChecksum(boxIDs)
	fmt.Println("Checksum:", checksum)

	firstID, secondID, index := findSimilarBoxIDs(boxIDs)
	fmt.Printf("Similar box IDs (differing index %d):\n\t%s\n\t%s\n", index, firstID, secondID)

	common := firstID[:index] + secondID[index+1:]
	fmt.Printf("Common letters:\n\t%s\n", common)
}
