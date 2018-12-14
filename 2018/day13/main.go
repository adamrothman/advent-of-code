package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Direction byte

const (
	DirectionUp    = 94  // caret "^"
	DirectionDown  = 118 // letter v "v"
	DirectionLeft  = 60  // less than "<"
	DirectionRight = 62  // greater than ">"
)

func (d Direction) String() string {
	return fmt.Sprintf("%c", d)
}

type TrackType byte

const (
	TrackTypeNone          = 32  // space " "
	TrackTypeVertical      = 124 // pipe "|"
	TrackTypeHorizontal    = 45  // dash "-"
	TrackTypeCurveForward  = 47  // forward slash "/"
	TrackTypeCurveBackward = 92  // backslash "\"
	TrackTypeIntersection  = 43  // plus "+"
)

func (t TrackType) String() string {
	return fmt.Sprintf("%c", t)
}

type TurnDirection int

const (
	TurnDirectionLeft     = iota
	TurnDirectionStraight = iota
	TurnDirectionRight    = iota
)

func (td TurnDirection) String() string {
	switch td {
	case TurnDirectionLeft:
		return "left"
	case TurnDirectionStraight:
		return "straight"
	case TurnDirectionRight:
		return "right"
	default:
		return ""
	}
}

type Point struct {
	X, Y int
}

type Cart struct {
	ID               int
	CurrentDirection Direction
	NextTurn         TurnDirection
}

func (c *Cart) UpdateDirection(t TrackType) {
	switch t {
	case TrackTypeCurveForward, TrackTypeCurveBackward:
		c.followCurve(t)
	case TrackTypeIntersection:
		c.turnAtIntersection()
	}
}

func (c *Cart) followCurve(t TrackType) {
	switch t {
	case TrackTypeCurveForward:
		switch c.CurrentDirection {
		case DirectionUp:
			c.CurrentDirection = DirectionRight
		case DirectionDown:
			c.CurrentDirection = DirectionLeft
		case DirectionLeft:
			c.CurrentDirection = DirectionDown
		case DirectionRight:
			c.CurrentDirection = DirectionUp
		}
	case TrackTypeCurveBackward:
		switch c.CurrentDirection {
		case DirectionUp:
			c.CurrentDirection = DirectionLeft
		case DirectionDown:
			c.CurrentDirection = DirectionRight
		case DirectionLeft:
			c.CurrentDirection = DirectionUp
		case DirectionRight:
			c.CurrentDirection = DirectionDown
		}
	}
}

func (c *Cart) turnAtIntersection() {
	switch c.NextTurn {
	case TurnDirectionLeft:
		switch c.CurrentDirection {
		case DirectionUp:
			c.CurrentDirection = DirectionLeft
		case DirectionDown:
			c.CurrentDirection = DirectionRight
		case DirectionLeft:
			c.CurrentDirection = DirectionDown
		case DirectionRight:
			c.CurrentDirection = DirectionUp
		}

		c.NextTurn = TurnDirectionStraight
	case TurnDirectionStraight:
		c.NextTurn = TurnDirectionRight
	case TurnDirectionRight:
		switch c.CurrentDirection {
		case DirectionUp:
			c.CurrentDirection = DirectionRight
		case DirectionDown:
			c.CurrentDirection = DirectionLeft
		case DirectionLeft:
			c.CurrentDirection = DirectionUp
		case DirectionRight:
			c.CurrentDirection = DirectionDown
		}

		c.NextTurn = TurnDirectionLeft
	}
}

type World struct {
	Width, Height int
	Tracks        map[Point]TrackType
	Carts         map[Point]*Cart
}

func (w World) String() string {
	builder := strings.Builder{}
	for y := 0; y < w.Height; y++ {
		lineBuilder := strings.Builder{}
		for x := 0; x < w.Width; x++ {
			p := Point{X: x, Y: y}
			if cart, ok := w.Carts[p]; ok {
				lineBuilder.WriteByte(byte(cart.CurrentDirection))
			} else {
				lineBuilder.WriteByte(byte(w.Tracks[p]))
			}
		}
		lineBuilder.WriteString("\n")
		builder.WriteString(lineBuilder.String())
	}
	return builder.String()
}

func readInput(filename string) (World, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return World{}, fmt.Errorf("constructing absolute path from %s: %s", filename, err)
	}

	f, err := os.Open(path)
	if err != nil {
		return World{}, fmt.Errorf("opening input file %s: %s", path, err)
	}
	defer f.Close()

	maxX, maxY, cartID := 0, 0, 0
	tracks := make(map[Point]TrackType)
	carts := make(map[Point]*Cart)

	scanner := bufio.NewScanner(f)
	for y := 0; scanner.Scan(); y++ {
		line := scanner.Text()

		for x := 0; x < len(line); x++ {
			b := line[x]
			p := Point{X: x, Y: y}

			if b == DirectionUp || b == DirectionDown || b == DirectionLeft || b == DirectionRight {
				// It's a cart, on a piece of track
				if b == DirectionUp || b == DirectionDown {
					tracks[p] = TrackTypeVertical
				} else {
					tracks[p] = TrackTypeHorizontal
				}

				carts[p] = &Cart{
					ID:               cartID,
					CurrentDirection: Direction(b),
					NextTurn:         TurnDirectionLeft,
				}

				cartID++
			} else {
				// It's a piece of track (or empty)
				tracks[p] = TrackType(b)
			}

			if x > maxX {
				maxX = x
			}
		}

		if y > maxY {
			maxY = y
		}
	}
	if err := scanner.Err(); err != nil {
		return World{}, fmt.Errorf("reading input file: %s", err)
	}

	world := World{
		Width:  maxX + 1,
		Height: maxY + 1,
		Tracks: tracks,
		Carts:  carts,
	}
	return world, nil
}

func simulate(world *World) {
	for tick := 0; len(world.Carts) > 1; tick++ {
		updatedCarts := make(map[int]bool)

		for y := 0; y < world.Height; y++ {
			for x := 0; x < world.Width; x++ {
				p := Point{X: x, Y: y}
				cart, ok := world.Carts[p]
				if !ok || updatedCarts[cart.ID] {
					continue
				}

				updatedCarts[cart.ID] = true

				var nextP Point
				switch cart.CurrentDirection {
				case DirectionUp:
					nextP = Point{X: p.X, Y: p.Y - 1}
				case DirectionDown:
					nextP = Point{X: p.X, Y: p.Y + 1}
				case DirectionLeft:
					nextP = Point{X: p.X - 1, Y: p.Y}
				case DirectionRight:
					nextP = Point{X: p.X + 1, Y: p.Y}
				}

				// If there's already a cart at nextP, we've got a crash!
				if _, ok := world.Carts[nextP]; ok {
					fmt.Printf("[%d]\tCrash at %+v\n", tick, nextP)

					delete(world.Carts, p)     // remove moving cart
					delete(world.Carts, nextP) // remove crashing cart

					continue
				}

				// Depending on the next track section, we may need to update
				// the cart's direction. We don't need to worry about going
				// off the edge of the world, or whether next track is valid
				// given current track.
				nextT := world.Tracks[nextP]
				cart.UpdateDirection(nextT)

				// Update the cart's location
				delete(world.Carts, p)
				world.Carts[nextP] = cart
			}
		}
	}
}

func main() {
	filename := "input.txt"

	world, err := readInput(filename)
	if err != nil {
		log.Printf("Error reading input from %s: %s\n", filename, err)
	}

	fmt.Println("Initial state")
	fmt.Println(world)

	simulate(&world)
	for p := range world.Carts {
		fmt.Printf("\nLast cart standing is at %+v\n", p)
	}
}
