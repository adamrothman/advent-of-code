package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

type Point struct {
	X int64
	Y int64
}

var pointRegexp = regexp.MustCompile(`^(\d+), (\d+)$`)

func parsePoint(raw string) (Point, error) {
	matches := pointRegexp.FindStringSubmatch(raw)
	if matches == nil {
		return Point{}, fmt.Errorf("point string \"%s\" does not match pattern", raw)
	}

	x, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return Point{}, fmt.Errorf("parsing X coordinate: %s", err)
	}

	y, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return Point{}, fmt.Errorf("parsing Y coordinate: %s", err)
	}

	return Point{X: x, Y: y}, nil
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

	scanner := bufio.NewScanner(f)
	points := make([]Point, 0)

	for scanner.Scan() {
		point, err := parsePoint(scanner.Text())
		if err != nil {
			log.Printf("Error parsing line: %s", err)
			continue
		}
		points = append(points, point)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input file: %s", err)
	}

	return points, nil
}

func manhattanDistance(p, q Point) int64 {
	return abs64(p.X-q.X) + abs64(p.Y-q.Y)
}

func abs64(x int64) int64 {
	y := x >> 63
	return (x ^ y) - y
}

type Bounds struct {
	MinX int64
	MinY int64
	MaxX int64
	MaxY int64
}

func findBounds(points []Point) Bounds {
	bounds := Bounds{
		MinX: math.MaxInt64,
		MinY: math.MaxInt64,
		MaxX: math.MinInt64,
		MaxY: math.MinInt64,
	}

	for _, point := range points {
		if point.X < bounds.MinX {
			bounds.MinX = point.X
		}
		if point.X > bounds.MaxX {
			bounds.MaxX = point.X
		}
		if point.Y < bounds.MinY {
			bounds.MinY = point.Y
		}
		if point.Y > bounds.MaxY {
			bounds.MaxY = point.Y
		}
	}

	return bounds
}

func findMostIsolatedPoint(points []Point) (Point, uint64) {
	areaByPoint := make(map[Point]uint64)
	infinite := make(map[Point]bool)

	bounds := findBounds(points)

	for x := bounds.MinX; x <= bounds.MaxX; x++ {
		for y := bounds.MinY; y <= bounds.MaxY; y++ {
			current := Point{X: x, Y: y}

			var smallestDistance int64 = math.MaxInt64
			var closestPoint *Point

			for _, point := range points {
				distance := manhattanDistance(point, current)

				if distance < smallestDistance {
					smallestDistance = distance
					closestPoint = &Point{X: point.X, Y: point.Y}
				} else if distance == smallestDistance {
					// Tie between 2+ points
					closestPoint = nil
				}
			}

			if closestPoint == nil {
				continue
			}

			if x == bounds.MinX || x == bounds.MaxX || y == bounds.MinY || y == bounds.MaxY {
				infinite[*closestPoint] = true
			} else {
				areaByPoint[*closestPoint]++
			}
		}
	}

	var mostIsolatedPoint Point
	var maxArea uint64
	for point, area := range areaByPoint {
		if infinite[point] {
			continue
		}
		if area > maxArea {
			maxArea = area
			mostIsolatedPoint = point
		}
	}

	return mostIsolatedPoint, maxArea
}

func findSafestRegionArea(points []Point) (area uint64) {
	seen := make(map[Point]bool)

	queue := make([]Point, 0)
	// Starting point doesn't matter, just pick the first of the ones provided
	queue = append(queue, points[0])

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if seen[current] {
			continue
		}
		seen[current] = true

		if totalManhattanDistance(current, points) < 10000 {
			area++
			queue = append(
				queue,
				Point{X: current.X + 1, Y: current.Y},
				Point{X: current.X - 1, Y: current.Y},
				Point{X: current.X, Y: current.Y + 1},
				Point{X: current.X, Y: current.Y - 1},
			)
		}
	}

	return
}

func totalManhattanDistance(ref Point, points []Point) (distance int64) {
	for _, p := range points {
		distance += manhattanDistance(ref, p)
	}
	return
}

func main() {
	filename := "input.txt"

	points, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	mostIsolatedPoint, area := findMostIsolatedPoint(points)
	fmt.Printf("Most isolated point is %+v with area %d\n", mostIsolatedPoint, area)

	area = findSafestRegionArea(points)
	fmt.Printf("Area of region containing all locations with total distance < 10000: %d\n", area)
}
