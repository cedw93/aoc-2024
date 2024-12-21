package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

type location struct {
	row               int
	col               int
	value             int
	visited           bool
	visitedFromParent map[string]struct{}
	score             int
	seen              int
}

type directionMap struct {
	rowOffset int
	colOffset int
	label     string
}

var grid [][]*location
var trailHeads []*location

var directions = [...]directionMap{
	{0, 1, "right"},
	{1, 0, "down"},
	{0, -1, "left"},
	{-1, 0, "up"},
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	rowCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		result := []*location{}
		for i, char := range line {
			location := &location{
				row:               rowCount,
				col:               i,
				value:             aToIIgnoreError(string(char)),
				visited:           false,
				visitedFromParent: make(map[string]struct{}),
				score:             0,
				seen:              0,
			}
			if char == '0' {
				location.seen++
				trailHeads = append(trailHeads, location)
			}

			result = append(result, location)
		}
		grid = append(grid, result)
		rowCount++
	}
}

func onGrid(row, col int, g [][]*location) bool {
	return row > -1 && col > -1 && row < len(g) && col < len(g[row])
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func (l *location) getNext() []*location {
	result := []*location{}
	for _, direction := range directions {
		rowOffset := l.row + direction.rowOffset
		colOffset := l.col + direction.colOffset
		if onGrid(rowOffset, colOffset, grid) {
			candidate := grid[rowOffset][colOffset]
			if l.value+1 == candidate.value {
				// We add the previous nodes seen as this is the number of distinct ways to get to this node
				candidate.seen += l.seen
				if !candidate.visited {
					candidate.visited = true
					result = append(result, candidate)
				}
			}
		}
	}
	return result
}

func resetGrid(grid [][]*location) {
	for _, row := range grid {
		for _, loc := range row {
			loc.visited = false
			// 0s are head nodes so seen will always be 1 no need to reset
			if loc.value != 0 {
				loc.seen = 0
			}
		}
	}
}

func bothParts() (int, int) {
	partOne := 0
	partTwo := 0

	// simple BFS
	for _, head := range trailHeads {
		queue := []*location{head}
		count := 0
		for len(queue) > 0 {
			curr := queue[0]
			curr.visited = true
			queue = queue[1:]
			if curr.value == 9 {
				// Part one cares about the score of the 'head'
				head.score += 1
				// part two cares about distinct ways to reach values of '9'
				partTwo += curr.seen
				continue
			}
			queue = append(queue, curr.getNext()...)
			count++
		}

		partOne += head.score
		resetGrid(grid)
	}

	return partOne, partTwo
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func main() {
	defer timer()()
	partOne, partTwo := bothParts()
	fmt.Println("Part One:", partOne)
	fmt.Println("Part Two:", partTwo)
}
