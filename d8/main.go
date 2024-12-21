package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var antiNodes [][]rune

type coordinate struct {
	row int
	col int
}

var antennas = make(map[rune][]coordinate)

const (
	ANTINODE    = '#'
	IGNORE_RUNE = '.'
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		antiNodeLine := []rune{}
		for i, r := range line {
			if r != IGNORE_RUNE {
				if _, ok := antennas[r]; ok {
					antennas[r] = append(antennas[r], coordinate{rowCount, i})
				} else {
					antennas[r] = []coordinate{{rowCount, i}}
				}
			}
			antiNodeLine = append(antiNodeLine, '.')
		}
		antiNodes = append(antiNodes, antiNodeLine)
		rowCount++
	}
}

func onGrid(row, col int, g [][]rune) bool {
	return row > -1 && col > -1 && row < len(g) && col < len(g[row])
}

// not using math.abs to avoid 100s of type conversions to float64
func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func addAntiNodesForLine(c1, c2 coordinate, g [][]rune) {

	// since a pair of nodes is always on a line set these to anti nodes too!

	g[c1.row][c1.col] = ANTINODE
	g[c2.row][c2.col] = ANTINODE

	rowDelta := c1.row - c2.row
	colDelta := c1.col - c2.col
	targetRow := c1.row + rowDelta
	targetCol := c1.col + colDelta

	for onGrid(targetRow, targetCol, g) {
		g[targetRow][targetCol] = ANTINODE
		targetRow += rowDelta
		targetCol += colDelta
	}

	rowDelta *= -1
	colDelta *= -1
	targetRow = c2.row + rowDelta
	targetCol = c2.col + colDelta

	for onGrid(targetRow, targetCol, g) {
		g[targetRow][targetCol] = ANTINODE
		targetRow += rowDelta
		targetCol += colDelta
	}
}

func addAntiNodesForPair(c1, c2 coordinate, g [][]rune) {
	// diff between two points, ordering doesnt really matter
	rowDelta := c1.row - c2.row
	colDelta := c1.col - c2.col
	targetRow := c1.row + rowDelta
	targetCol := c1.col + colDelta

	// fmt.Printf("Delta is: (%d,%d)\n", rowDelta, colDelta)
	// fmt.Printf("Current (%d,%d). Attempting to add AntiNode at (%d, %d)\n", c1.row, c1.col, targetRow, targetCol)

	if onGrid(targetRow, targetCol, g) {
		g[targetRow][targetCol] = ANTINODE
	}

	// Invert them as we need the 'mirrored' copy
	rowDelta *= -1
	colDelta *= -1
	targetRow = c2.row + rowDelta
	targetCol = c2.col + colDelta

	// fmt.Printf("Delta is: (%d,%d)\n", rowDelta, colDelta)
	// fmt.Printf("Current (%d,%d). Attempting to add AntiNode at (%d, %d)\n", c2.row, c2.col, targetRow, targetCol)

	if onGrid(targetRow, targetCol, g) {
		g[targetRow][targetCol] = ANTINODE
	}
}

func partOne() int {
	result := 0
	for _, coordinates := range antennas {
		// iterate the pairs
		for i := 0; i < len(coordinates); i++ {
			// j = i + 1 otherwise were just checking itself against itself
			for j := i + 1; j < len(coordinates); j++ {
				addAntiNodesForPair(coordinates[i], coordinates[j], antiNodes)
			}
		}

	}

	for _, row := range antiNodes {
		for _, r := range row {
			if r == ANTINODE {
				result += 1
			}
		}
	}
	return result
}

func partTwo() int {
	result := 0

	for _, coordinates := range antennas {
		for i := 0; i < len(coordinates); i++ {
			for j := i + 1; j < len(coordinates); j++ {
				// Reused whats been done on step 1, no need to re do pairs
				addAntiNodesForLine(coordinates[i], coordinates[j], antiNodes)
			}
		}
	}

	for _, row := range antiNodes {
		for _, r := range row {
			if r == ANTINODE {
				result += 1
			}
		}
	}

	return result
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())
}
