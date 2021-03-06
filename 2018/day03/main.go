package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Claim struct {
	ID uint64

	Left uint64
	Top  uint64

	Width  uint64
	Height uint64
}

func (c Claim) Overlaps(other Claim) bool {
	if c.Left+c.Width < other.Left || other.Left+other.Width < c.Left {
		return false
	}
	if c.Top+c.Height < other.Top || other.Top+other.Height < c.Top {
		return false
	}
	return true
}

func readInput(filename string) ([]Claim, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	claims := make([]Claim, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var c Claim
		n, err := fmt.Sscanf(
			scanner.Text(),
			"#%d @ %d,%d: %dx%d",
			&c.ID,
			&c.Left,
			&c.Top,
			&c.Width,
			&c.Height,
		)
		if n != 5 || err != nil {
			return nil, fmt.Errorf("parsing claim: %s", err)
		}
		claims = append(claims, c)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return claims, nil
}

type Fabric [][]int

const FabricDimension = 1000

func NewFabric(dimension int) Fabric {
	fabric := make(Fabric, dimension)
	for i := 0; i < dimension; i++ {
		fabric[i] = make([]int, dimension)
	}
	return fabric
}

func populateFabric(claims []Claim) Fabric {
	fabric := NewFabric(FabricDimension)

	for _, claim := range claims {
		for x := claim.Left; x < claim.Left+claim.Width; x++ {
			for y := claim.Top; y < claim.Top+claim.Height; y++ {
				fabric[x][y]++
			}
		}
	}

	return fabric
}

func countOverlappingSquares(fabric Fabric) (count int) {
	for _, col := range fabric {
		for _, claimCount := range col {
			if claimCount > 1 {
				count++
			}
		}
	}

	return
}

func findNonOverlappingClaim(claims []Claim) uint64 {
	hasOverlap := make(map[uint64]bool)
	for _, claim := range claims {
		hasOverlap[claim.ID] = false
	}

	for _, candidate := range claims {
		if hasOverlap[candidate.ID] {
			continue
		}

		for _, other := range claims {
			if other.ID == candidate.ID {
				continue
			}

			if candidate.Overlaps(other) {
				hasOverlap[candidate.ID] = true
				hasOverlap[other.ID] = true
			}
		}
	}

	for id, overlaps := range hasOverlap {
		if !overlaps {
			return id
		}
	}

	return 0
}

func main() {
	filename := "input.txt"

	claims, err := readInput(filename)
	if err != nil {
		log.Fatalf("Error reading input from %s: %s\n", filename, err)
	}

	fabric := populateFabric(claims)

	overlapCount := countOverlappingSquares(fabric)
	fmt.Println("Square inches with overlap:", overlapCount)

	nonOverlappingID := findNonOverlappingClaim(claims)
	fmt.Println("Non-overlapping claim ID:", nonOverlappingID)
}
