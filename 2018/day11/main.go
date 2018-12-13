package main

import (
	"fmt"
	"math"
)

type Cell struct {
	X, Y int64
}

type PowerGrid struct {
	Width, Height int64
	Serial        int64

	summedPowerTable map[Cell]int64
}

func NewPowerGrid(width, height, serial int64) PowerGrid {
	grid := PowerGrid{
		Width:  width,
		Height: height,
		Serial: serial,
	}

	// https://en.wikipedia.org/wiki/Summed-area_table
	spt := make(map[Cell]int64)
	for x := int64(1); x <= width; x++ {
		for y := int64(1); y <= height; y++ {
			power := grid.CalculatePowerAt(x, y)
			spt[Cell{X: x, Y: y}] = power + spt[Cell{X: x, Y: y - 1}] + spt[Cell{X: x - 1, Y: y}] - spt[Cell{X: x - 1, Y: y - 1}]
		}
	}
	grid.summedPowerTable = spt

	return grid
}

func (g PowerGrid) CalculatePowerAt(x, y int64) (power int64) {
	rackID := x + 10
	power = rackID * y
	power += g.Serial
	power *= rackID
	power = (power / 100) % 10 // hundreds digit
	power -= 5
	return
}

// https://en.wikipedia.org/wiki/Summed-area_table#The_algorithm
// Note that our A here is actually 1 unit to the left and above the given
// topLeft; this is because the canonical formula excludes the top and left
// sides (which we want to include).
func (g PowerGrid) GetPowerWithin(topLeft, bottomRight Cell) int64 {
	a := Cell{X: topLeft.X - 1, Y: topLeft.Y - 1}
	b := Cell{X: bottomRight.X, Y: a.Y}
	c := Cell{X: a.X, Y: bottomRight.Y}
	d := Cell{X: bottomRight.X, Y: bottomRight.Y}
	return g.summedPowerTable[d] + g.summedPowerTable[a] - g.summedPowerTable[b] - g.summedPowerTable[c]
}

func findLargestTotalPower(grid PowerGrid, width, height int64) (topLeft Cell, totalPower int64) {
	totalPower = math.MinInt64

	for x := int64(1); x+width <= grid.Width; x++ {
		for y := int64(1); y+height <= grid.Height; y++ {
			tL := Cell{X: x, Y: y}
			bR := Cell{X: x + width - 1, Y: y + height - 1}
			if power := grid.GetPowerWithin(tL, bR); power > totalPower {
				topLeft = tL
				totalPower = power
			}
		}
	}

	return
}

func findOverallLargestTotalPower(grid PowerGrid) (topLeft Cell, bestSize, totalPower int64) {
	totalPower = math.MinInt64

	for size := int64(1); size < grid.Width; size++ {
		for x := int64(1); x+size <= grid.Width; x++ {
			for y := int64(1); y+size <= grid.Height; y++ {
				tL := Cell{X: x, Y: y}
				bR := Cell{X: x + size - 1, Y: y + size - 1}
				if power := grid.GetPowerWithin(tL, bR); power > totalPower {
					topLeft = tL
					bestSize = size
					totalPower = power
				}
			}
		}
	}

	return
}

func main() {
	var serial int64 = 4842

	grid := NewPowerGrid(300, 300, serial)

	topLeft, totalPower := findLargestTotalPower(grid, 3, 3)
	fmt.Printf("Top left cell of 3x3 square with largest total power: %+v (power %d)\n", topLeft, totalPower)

	topLeft, bestSize, totalPower := findOverallLargestTotalPower(grid)
	fmt.Printf("Square with largest total power has top left %+v and size %d (total power %d)\n", topLeft, bestSize, totalPower)
}
