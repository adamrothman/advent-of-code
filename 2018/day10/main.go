package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
)

type point struct {
	X, Y int64
}

type Point struct {
	Position point
	Velocity point
}

func readInput(filename string) ([]Point, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	points := make([]Point, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var p Point
		n, err := fmt.Sscanf(
			scanner.Text(),
			"position=<%d, %d> velocity=<%d, %d>",
			&p.Position.X,
			&p.Position.Y,
			&p.Velocity.X,
			&p.Velocity.Y,
		)
		if n != 4 || err != nil {
			return nil, fmt.Errorf("parsing line: %s", err)
		}
		points = append(points, p)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return points, nil
}

func findMessageArrangement(points []Point) (arrangement map[point]bool, messageT int64) {
	var minArea int64 = math.MaxInt64

	for t := int64(0); ; t++ {
		pointsAtT := plot(points, t)
		area := calculateArea(pointsAtT)
		if area < minArea {
			minArea = area
			arrangement = pointsAtT
			messageT = t
		} else {
			break
		}
	}

	return
}

func plot(points []Point, t int64) map[point]bool {
	result := make(map[point]bool)

	for _, p := range points {
		pAtT := point{
			X: p.Position.X + p.Velocity.X*t,
			Y: p.Position.Y + p.Velocity.Y*t,
		}
		result[pAtT] = true
	}

	return result
}

func calculateArea(points map[point]bool) int64 {
	min, max := findBounds(points)
	return (max.X - min.X) * (max.Y - min.Y)
}

func findBounds(points map[point]bool) (min, max point) {
	min.X, min.Y = math.MaxInt64, math.MaxInt64
	max.X, max.Y = math.MinInt64, math.MinInt64

	for p := range points {
		if p.X < min.X {
			min.X = p.X
		}
		if p.X > max.X {
			max.X = p.X
		}
		if p.Y < min.Y {
			min.Y = p.Y
		}
		if p.Y > max.Y {
			max.Y = p.Y
		}
	}

	return
}

func draw(points map[point]bool) {
	min, max := findBounds(points)
	for y := min.Y - 1; y <= max.Y+1; y++ {
		for x := min.X - 1; x <= max.X+1; x++ {
			if points[point{X: x, Y: y}] {
				fmt.Printf("#")
			} else {
				fmt.Printf(".")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func main() {
	filename := "input.txt"

	points, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	arrangement, t := findMessageArrangement(points)
	fmt.Printf("After %d seconds:\n", t)
	draw(arrangement)
}
