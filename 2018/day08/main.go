package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func readInput(filename string) ([]uint64, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	values := make([]uint64, 0)

	for scanner.Scan() {
		value, err := strconv.ParseUint(scanner.Text(), 10, 64)
		if err != nil {
			log.Printf("Error parsing uint64: %s", err)
			continue
		}
		values = append(values, value)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return values, nil
}

type Node struct {
	Children []Node
	Metadata []uint64
}

func buildTree(numbers []uint64, startIndex uint64) (node Node, consumedCount uint64) {
	childCount, metadataCount := numbers[startIndex], numbers[startIndex+1]
	consumedCount = 2

	node = Node{
		Children: make([]Node, 0, childCount),
		Metadata: make([]uint64, 0, metadataCount),
	}

	currentIndex := startIndex + consumedCount

	for c := uint64(0); c < childCount; c++ {
		child, consumed := buildTree(numbers, currentIndex)

		node.Children = append(node.Children, child)

		consumedCount += consumed
		currentIndex += consumed
	}

	for m := uint64(0); m < metadataCount; m++ {
		entry := numbers[currentIndex+m]
		node.Metadata = append(node.Metadata, entry)
	}

	consumedCount += metadataCount

	return
}

func sumMetadata(root Node) (metadataSum uint64) {
	for _, child := range root.Children {
		metadataSum += sumMetadata(child)
	}

	for _, entry := range root.Metadata {
		metadataSum += entry
	}

	return
}

func computeValue(node Node) (value uint64) {
	if len(node.Children) == 0 {
		for _, entry := range node.Metadata {
			value += entry
		}
		return
	}

	for _, entry := range node.Metadata {
		if entry == 0 {
			continue
		}

		index := int(entry) - 1
		if index >= len(node.Children) {
			continue
		}

		child := node.Children[index]
		value += computeValue(child)
	}

	return
}

func main() {
	filename := "input.txt"

	numbers, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	tree, _ := buildTree(numbers, 0)
	sum := sumMetadata(tree)
	fmt.Println("Sum of metadata entries:", sum)

	value := computeValue(tree)
	fmt.Println("Tree value:", value)
}
