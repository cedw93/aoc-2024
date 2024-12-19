package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type direction int

type guard struct {
	row                 int
	col                 int
	direction           direction
	onGrid              bool
	numPositionsVisited int
	// part one, tracks if a node is distinct or not
	visited [][]bool
	// part two, tracks if weve been to a node whilst facing a certain direction before
	visitedWithDirection map[direction]map[string]struct{}
	inLoop               bool
}

var grid [][]rune

const (
	GUARD_START_RUNE           = '^'
	OBSTACLE_RUNE              = '#'
	SAFE_SPACE                 = '.'
	UP               direction = 0
	DOWN             direction = 1
	LEFT             direction = 2
	RIGHT            direction = 3
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	count := 0
	for scanner.Scan() {
		grid = append(grid, []rune(scanner.Text()))
		count++
	}
}

// Simple clockwise mappings, could likely make this an actual map if it got much bigger
func (g *guard) updateDirection() {
	if g.direction == UP {
		g.direction = RIGHT
	} else if g.direction == DOWN {
		g.direction = LEFT
	} else if g.direction == LEFT {
		g.direction = UP
	} else if g.direction == RIGHT {
		g.direction = DOWN
	}

}

// Each move is exactly one tile
func (g *guard) nextMoveLocation() (int, int) {
	offsetMatrix := [2]int{0, 0}
	if g.direction == UP {
		offsetMatrix[0], offsetMatrix[1] = -1, 0
	} else if g.direction == DOWN {
		offsetMatrix[0], offsetMatrix[1] = 1, 0
	} else if g.direction == LEFT {
		offsetMatrix[0], offsetMatrix[1] = 0, -1
	} else if g.direction == RIGHT {
		offsetMatrix[0], offsetMatrix[1] = 0, 1
	}

	return g.row + offsetMatrix[0], g.col + offsetMatrix[1]
}

// moves the guard, checks are done before this function is called
// so this is ALWAYS safe.
// also tracks if this move is a new node, or if this move will cause a loop
func (g *guard) move(row, col int) {
	g.row = row
	g.col = col
	directionData := g.visitedWithDirection[g.direction]
	targetKey := fmt.Sprintf("%d-%d", g.row, g.col)
	if visited := g.visited[row][col]; !visited {
		g.visited[row][col] = true
		g.numPositionsVisited++
	}

	// targetKey will be in the form rol-col
	// so we look in the map[direction] if key row-col exists if it does we've been to [row][col] before whilst facing this direction
	// which will cause a loop. Its fine to come back to this tile if our direction is different
	if _, ok := directionData[targetKey]; ok {
		g.inLoop = true
		g.onGrid = false
	} else {
		directionData[targetKey] = struct{}{}
	}
}

func (g *guard) printVisited() {
	for _, row := range g.visited {
		for _, state := range row {
			fmt.Printf("%t,", state)
		}
		fmt.Println()
	}
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func findGuard(g [][]rune) *guard {

	dirMap := make(map[direction]map[string]struct{})
	dirMap[UP] = make(map[string]struct{})
	dirMap[DOWN] = make(map[string]struct{})
	dirMap[LEFT] = make(map[string]struct{})
	dirMap[RIGHT] = make(map[string]struct{})

	result := &guard{
		row:                  -1,
		col:                  -1,
		direction:            UP,
		onGrid:               false,
		numPositionsVisited:  0,
		visited:              [][]bool{},
		visitedWithDirection: dirMap,
		inLoop:               false,
	}

	for rIdx, row := range g {
		visitedRow := make([]bool, len(row), len(row))
		for cIdx, r := range row {
			visitedRow[cIdx] = false
			if r == GUARD_START_RUNE {
				result.row = rIdx
				result.col = cIdx
				result.onGrid = true
				// track starting location
				visitedRow[cIdx] = true
				result.numPositionsVisited++
			}
		}

		result.visited = append(result.visited, visitedRow)
	}

	return result

}

func onGrid(row, col int, g [][]rune) bool {
	return row > -1 && col > -1 && row < len(g) && col < len(g[row])
}

func partOne() ([][]bool, int) {
	g := findGuard(grid)
	// were still on the board and we are not looping
	for g.onGrid && !g.inLoop {
		// don't actually move just based on direction facing get the next row, col indexes
		// This could be way more efficient as one we are in a direction we can skip ahead until end or object
		// you don't need to check every tile, but this is the most simple approach with less edge cases and input is small enough for it not to be that bad
		nextRow, nextCol := g.nextMoveLocation()

		// is the move outside the grid?
		if !onGrid(nextRow, nextCol, grid) {
			g.onGrid = false
			continue
		}

		nextRune := grid[nextRow][nextCol]

		// is the next move an obstacle? if so, rotate 90 and try again. Guard will still be on same row, col just rotated
		if nextRune == OBSTACLE_RUNE {
			g.updateDirection()
			continue
		}

		// Actually move, this is safe as conditions before it check bounds etc
		g.move(nextRow, nextCol)
	}

	return g.visited, g.numPositionsVisited
}

func partTwo(partOneVisited [][]bool) int {
	result := 0

	// for every [row][col] in the grid, replace with an obstacle if it's currently a safe space
	// then have the guard perform its routing on that new grid pattern
	// guards detect if they are in a loop already, otherwise same rules as part
	for rIdx, row := range grid {
		for cIdx := range row {
			// We only need check tiles that we know are on the valid path which was found in part 1
			// otherwise no point adding an obstacle somewhere the guard never walks!
			if partOneVisited[rIdx][cIdx] && grid[rIdx][cIdx] == SAFE_SPACE {
				grid[rIdx][cIdx] = OBSTACLE_RUNE
				g := findGuard(grid)
				for g.onGrid {
					nextRow, nextCol := g.nextMoveLocation()

					if !onGrid(nextRow, nextCol, grid) {
						g.onGrid = false
						continue
					}

					nextRune := grid[nextRow][nextCol]

					if nextRune == OBSTACLE_RUNE {
						g.updateDirection()
						continue
					}

					g.move(nextRow, nextCol)

					if g.inLoop {
						result++
						break
					}
				}
				grid[rIdx][cIdx] = SAFE_SPACE
			}
		}
	}
	return result
}

func main() {
	defer timer()()
	visited, distinctNodes := partOne()
	fmt.Println("Part One:", distinctNodes)
	fmt.Println("Part Two:", partTwo(visited))
}
