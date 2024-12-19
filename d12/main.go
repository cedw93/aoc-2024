package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var grid [][]*location
var regions []region

type location struct {
	row       int
	col       int
	value     rune
	perimeter int
	visited   bool
	corners   int
}

type region struct {
	value     string
	locations []*location
	perimeter int
	area      int
	price     int
	edges     int
}

type directionMap struct {
	rowOffset int
	colOffset int
	label     string
}

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
		row := []*location{}
		for i, r := range line {
			row = append(row, &location{
				row:       rowCount,
				col:       i,
				value:     r,
				perimeter: 0,
				visited:   false,
				corners:   0,
			})
		}
		grid = append(grid, row)
		rowCount++
	}
}

func (l *location) calcCorners() {
	// top left edge of grid
	if l.row == 0 && l.col == 0 {
		l.corners++
	}

	// top right edge of grid
	if l.row == 0 && l.col == len(grid)-1 {
		l.corners++
	}

	// bottom left edge of grid
	if l.row == len(grid[0])-1 && l.col == len(grid)-1 {
		l.corners++
	}

	// bottom right edge of grid
	if l.row == len(grid[0])-1 && l.col == 0 {
		l.corners++
	}

	// Now the messy corner checking... lets take [2,2] as an example for all the below
	//
	//	[1,1][1,2][1.2]
	//	[2,1][2,2][2,3]
	//	[3.1][3,2][3,3]
	//

	// Top left outside corner
	// this would check (and boundary checks)
	// ([2,1] != l.value && [1,2] != l.value ) || ([2,1] != l.value) || ([1,2] != l.value)
	// means that in [2, 2] is an outside left corner
	// and outside corner is an outside corner if X and Y != A where A is the regionID
	// .Y.
	// XAA
	// .AA
	if (l.col > 0 && l.row > 0 && grid[l.row][l.col-1].value != l.value && grid[l.row-1][l.col].value != l.value) || (l.col > 0 && l.row == 0 && grid[l.row][l.col-1].value != l.value) || (l.col == 0 && l.row > 0 && grid[l.row-1][l.col].value != l.value) {
		l.corners++
	}

	// Top left inside corner
	// this would check (and boundary checks)
	// [2,3] == l.value && [3.2] == l.value && [3,3] != l.value
	// means that in [2, 2] is an inside corner
	// and inside corner is an inside corner if X != A where A is the regionID
	// AAA
	// AAA
	// AAX
	if l.col < len(grid[0])-1 && l.row < len(grid)-1 && grid[l.row][l.col+1].value == l.value && grid[l.row+1][l.col].value == l.value && grid[l.row+1][l.col+1].value != l.value {
		l.corners++
	}

	//bottom right outside
	if (l.col < len(grid[0])-1 && l.row < len(grid)-1 && grid[l.row][l.col+1].value != l.value && grid[l.row+1][l.col].value != l.value) || (l.col < len(grid[0])-1 && l.row == len(grid)-1 && grid[l.row][l.col+1].value != l.value) || (l.col == len(grid[0])-1 && l.row < len(grid)-1 && grid[l.row+1][l.col].value != l.value) {
		l.corners++
	}

	//bottom right inside
	if l.col > 0 && l.row > 0 && grid[l.row][l.col-1].value == l.value && grid[l.row-1][l.col].value == l.value && grid[l.row-1][l.col-1].value != l.value {
		l.corners++
	}

	// bottom left outside
	if (l.col > 0 && l.row < len(grid)-1 && grid[l.row][l.col-1].value != l.value && grid[l.row+1][l.col].value != l.value) || (l.col > 0 && l.row == len(grid)-1 && grid[l.row][l.col-1].value != l.value) || (l.col == 0 && l.row < len(grid)-1 && grid[l.row+1][l.col].value != l.value) {
		l.corners++
	}

	// bottom left inside
	if l.col < len(grid[0])-1 && l.row > 0 && grid[l.row][l.col+1].value == l.value && grid[l.row-1][l.col].value == l.value && grid[l.row-1][l.col+1].value != l.value {
		l.corners++
	}

	// top right outside
	if (l.col < len(grid[0])-1 && l.row > 0 && grid[l.row][l.col+1].value != l.value && grid[l.row-1][l.col].value != l.value) || (l.col < len(grid[0])-1 && l.row == 0 && grid[l.row][l.col+1].value != l.value) || (l.col == len(grid[0])-1 && l.row > 0 && grid[l.row-1][l.col].value != l.value) {
		l.corners++
	}

	//top right inside
	if l.col > 0 && l.row < len(grid)-1 && grid[l.row][l.col-1].value == l.value && grid[l.row+1][l.col].value == l.value && grid[l.row+1][l.col-1].value != l.value {
		l.corners++
	}

}

func (l *location) getNext() []*location {
	result := []*location{}
	for _, direction := range directions {
		rowOffset := l.row + direction.rowOffset
		colOffset := l.col + direction.colOffset
		if onGrid(rowOffset, colOffset, grid) {
			candidate := grid[rowOffset][colOffset]
			if l.value == candidate.value {
				if !candidate.visited {
					candidate.visited = true
					result = append(result, candidate)
				}
			} else {
				// candidate is not our region so this also must be an edge
				l.perimeter++
			}
		} else {
			// If we are not onGrid then current loc must be an edge
			l.perimeter++
		}
	}
	return result
}

func onGrid(row, col int, g [][]*location) bool {
	return row > -1 && col > -1 && row < len(g) && col < len(g[row])
}

func partOne() int {
	result := 0

	for _, row := range grid {
		for _, loc := range row {
			if loc.visited {
				continue
			}

			queue := []*location{loc}

			currentRegion := region{
				value:     string(loc.value),
				locations: []*location{},
				perimeter: 0,
				area:      0,
				price:     0,
				edges:     0,
			}

			for len(queue) > 0 {
				curr := queue[0]
				curr.visited = true
				currentRegion.locations = append(currentRegion.locations, curr)
				queue = queue[1:]
				queue = append(queue, curr.getNext()...)
				currentRegion.perimeter += curr.perimeter
			}

			currentRegion.area = len(currentRegion.locations)
			currentRegion.price = currentRegion.area * currentRegion.perimeter
			regions = append(regions, currentRegion)
			result += currentRegion.price
		}
	}

	return result
}

func partTwo() int {
	result := 0
	for _, region := range regions {
		for _, loc := range region.locations {
			loc.calcCorners()
			region.edges += loc.corners
		}
		result += region.area * region.edges
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
