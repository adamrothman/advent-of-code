package main

import (
	"fmt"
	"math"
)

type Cell struct {
	X, Y int64
}

func calculatePower(cell Cell, serial int64) (power int64) {
	rackID := cell.X + 10
	power = rackID * cell.Y
	power += serial
	power *= rackID
	power = (power / 100) % 10 // hundreds digit
	power -= 5
	return
}

type PowerGrid struct {
	Width, Height int64
	Serial        int64

	cellPowers map[Cell]int64
}

func NewPowerGrid(width, height, serial int64) PowerGrid {
	grid := PowerGrid{
		Width:      width,
		Height:     height,
		Serial:     serial,
		cellPowers: make(map[Cell]int64),
	}

	for x := int64(1); x <= width; x++ {
		for y := int64(1); y <= height; y++ {
			c := Cell{X: x, Y: y}
			grid.cellPowers[c] = calculatePower(c, serial)
		}
	}

	return grid
}

func (g PowerGrid) GetPowerAt(x, y int64) int64 {
	return g.cellPowers[Cell{X: x, Y: y}]
}

func (g PowerGrid) GetPowerWithin(topLeft, bottomRight Cell) (power int64) {
	for x := topLeft.X; x <= bottomRight.X; x++ {
		for y := topLeft.Y; y <= bottomRight.Y; y++ {
			power += g.GetPowerAt(x, y)
		}
	}
	return
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

// https://en.wikipedia.org/wiki/Summed-area_table
type SummedPowerTable struct {
	table map[Cell]int64
}

func NewSummedPowerTable(grid PowerGrid) SummedPowerTable {
	table := make(map[Cell]int64)

	for x := int64(1); x <= grid.Width; x++ {
		for y := int64(1); y <= grid.Height; y++ {
			table[Cell{X: x, Y: y}] = grid.GetPowerAt(x, y) + table[Cell{X: x, Y: y - 1}] + table[Cell{X: x - 1, Y: y}] - table[Cell{X: x - 1, Y: y - 1}]
		}
	}

	return SummedPowerTable{table: table}
}

// https://en.wikipedia.org/wiki/Summed-area_table#The_algorithm
// Note that our A here is actually 1 unit to the left and above the given
// topLeft; this is because the canonical formula excludes the top and left
// sides (which we want to include).
func (spt SummedPowerTable) GetPowerWithin(topLeft, bottomRight Cell) int64 {
	a := Cell{X: topLeft.X - 1, Y: topLeft.Y - 1}
	b := Cell{X: bottomRight.X, Y: a.Y}
	c := Cell{X: a.X, Y: bottomRight.Y}
	d := Cell{X: bottomRight.X, Y: bottomRight.Y}
	return spt.table[d] + spt.table[a] - spt.table[b] - spt.table[c]
}

func findOverallLargestTotalPower(grid PowerGrid) (topLeft Cell, bestSize, totalPower int64) {
	totalPower = math.MinInt64

	spt := NewSummedPowerTable(grid)

	for size := int64(1); size < grid.Width; size++ {
		for x := int64(1); x+size <= grid.Width; x++ {
			for y := int64(1); y+size <= grid.Height; y++ {
				tL := Cell{X: x, Y: y}
				bR := Cell{X: x + size - 1, Y: y + size - 1}
				if power := spt.GetPowerWithin(tL, bR); power > totalPower {
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
