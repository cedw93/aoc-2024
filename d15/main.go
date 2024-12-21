package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	WALL             = "#"
	ROBOT            = "@"
	EMPTY            = "."
	BOX              = "O"
	DOUBLE_BOX_LEFT  = "["
	DOUBLE_BOX_RIGHT = "]"
)

var gridInput []string
var instructions string

type grid [][]string

type direction struct {
	rowOffset int
	colOffset int
	value     rune
}

var directionMap = map[rune]direction{
	'>': {0, 1, '>'},
	'v': {1, 0, 'v'},
	'<': {0, -1, '<'},
	'^': {-1, 0, '^'},
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	var instructionBuff bytes.Buffer
	gridComplete := false
	for scanner.Scan() {
		line := scanner.Text()
		if !gridComplete {
			if len(strings.TrimSpace(line)) == 0 {
				gridComplete = true
				continue
			}
			gridInput = append(gridInput, line)
		} else {
			instructionBuff.WriteString(line)
		}
	}
	instructions = instructionBuff.String()
}

func createdGrid(lines []string) grid {
	var g grid
	for _, line := range lines {
		row := []string{}
		for _, r := range line {
			row = append(row, string(r))
		}
		g = append(g, row)

	}
	return g
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func (g grid) robotLocation() (int, int) {
	for r, row := range g {
		for c, v := range row {
			if v == ROBOT {
				return r, c
			}
		}
	}
	return -1, -1
}

func (g grid) print() {
	for _, row := range g {
		for _, v := range row {
			fmt.Printf("%s", string(v))
		}
		fmt.Println()
	}
}

// The GPS coordinate of a box is equal to 100 times its distance from the top edge of the map plus its distance from the left edge of the map
// Part 2 only cares about the nearest edge (so left box)
func (g grid) score() int {
	result := 0
	for r, row := range g {
		if r == 0 || r == len(g)-1 {
			continue
		}
		for c, value := range row {
			if value == BOX || value == DOUBLE_BOX_LEFT {
				result += (100 * r) + c
			}
		}
	}

	return result
}

func (g grid) onGrid(row, col int) bool {
	return row > -1 && col > -1 && row < len(g) && col < len(g[row])
}

// dfs here is used for horizontal movements only and boxes can'y overlap
// basically it continually moves every [ or ] if it can
// but only once it knows it safe to move
// recursively runs until it finds an empty (.) then that will return true
// then previous calls will move themselves along
// Robot is moved elsewhere, if this returns true then move the robot
// This could probably be a for-loop as its not really proper DFS
// this basically will do the follow given
// #.[][]@ -> recursive until it finds .
// #.[][]@ -> #[.][]@ (empty founds, swaps then returns true)
// #[.][]@ -> #[].[]@ recursive calls now see the true returns and start swapping
// #[].[]@ -> #[][.]@
// #[][.]@ -> #[][].@ back at our starting point, all funcs return, return true
// move robot separately if DFS returned true
// #[][].@ -> #[][]@.
// We've now moved everything left one tile!
func dfs(direction direction, row, col int, grid [][]string) bool {
	nextCol, nextRow := row+direction.rowOffset, col+direction.colOffset

	// Wall means nothing will move at all as there was no gap in the chain
	if grid[nextCol][nextRow] == WALL {
		return false
	}

	// We found a gap in the chain so simply swap the values in the grid around and return true
	// This should only happen once per chain
	if grid[nextCol][nextRow] == EMPTY {
		grid[nextCol][nextRow], grid[row][col] = grid[row][col], grid[nextCol][nextRow]
		return true
	}

	// If we are still checking boxes, then recursive DFS to see if movement if possible
	// if it is, then swap the adjacent grid cells based on the direction
	if grid[nextCol][nextRow] == DOUBLE_BOX_RIGHT || grid[nextCol][nextRow] == DOUBLE_BOX_LEFT {
		if !dfs(direction, nextCol, nextRow, grid) {
			return false
		}
		grid[row][col], grid[nextCol][nextRow] = grid[nextCol][nextRow], grid[row][col]
	}
	return true
}

// BFS for vertical as boxes can overlap
// will only allow movement if both spaces above or below all target boxes are empty (+ any new target boxes it finds)
// if any of them are not it will return false and won't move anything
// This is a bit more complication to demo
// but given
// ##############
// ##......##..##
// ##..........##
// ##...[][]...##
// ##....[]....##
// ##.....@....##
// ##############
// we attempt to move up with so start the search robot (@) [y-1, x]s
// queue = [[y-1, x]]
// g[y-1][x] == ] we add the index of [ to the queue
// queue = [[y-1, x], [y, x-1]]
// next = [y-1, x] (queue: [[y, x-1]])
// .. skip visited stuff ...
// track what [y, x] we are on so we know what to move later if needed
// neighbour := [next.Y-1, next.X] (-1 since we are moving up)
// if g[neighbour.Y][neighbour.X] == WALL -> return false, not allowed to move
// if g[neighbour.Y][neighbour.X] == EMPTY -> move onto next as the title has no importance yet
// if g[neighbour.Y][neighbour.X] == ONE_SIDE_OF_BOX, then add the other side to the queue and process
// Keep repeating, and if there is always space above the top most boxes (and space for all boxes) and move will happen
// otherwise if you hit a wall for any box, then none of them will be moved as it exits out
// Moves are only done at the end here!
func bfs(d direction, row, col int, g grid) bool {
	queue := [][]int{{
		row, col,
	}}
	visitedMap := make(map[string]struct{})
	visited := [][]int{}

	// Which side of the box are we?
	if g[row][col] == DOUBLE_BOX_RIGHT {
		queue = append(queue, []int{row, col - 1})
	} else {
		queue = append(queue, []int{row, col + 1})
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		row := curr[0]
		col := curr[1]
		targetKey := fmt.Sprintf("%d-%d", row, col)

		if _, ok := visitedMap[targetKey]; ok {
			// We've already visited here!
			continue
		}

		visitedMap[targetKey] = struct{}{}
		visited = append(visited, curr)

		nextRow, nextCol := row+d.rowOffset, col+d.colOffset

		if g[nextRow][nextCol] == WALL {
			// Can't move anything...
			return false
		}

		if g[nextRow][nextCol] == EMPTY {
			// nothing special to do here we will move at the end and keep searching
			continue
		}

		// If the next node is another part of a box add the box above or below to the list and we carry on!
		if g[nextRow][nextCol] == DOUBLE_BOX_LEFT {
			queue = append(queue, []int{nextRow, nextCol})
			queue = append(queue, []int{nextRow, nextCol + 1})
		}

		if g[nextRow][nextCol] == DOUBLE_BOX_RIGHT {
			queue = append(queue, []int{nextRow, nextCol})
			queue = append(queue, []int{nextRow, nextCol - 1})
		}
	}

	// BFS stored everything it visited now move them in reverse order to prevent weird swapping stuff, essentially we move
	// the empty spaces not the boxes.
	// variables are a bit messed up here, should have used a struct with proper names. horrible debugging statements
	for i := len(visited) - 1; i >= 0; i-- {
		visitedLocation := visited[i]
		r, c := visitedLocation[1]+d.colOffset, visitedLocation[0]+d.rowOffset
		// fmt.Println(visitedLocation[1], visitedLocation[0])
		// fmt.Println(r, c)
		// fmt.Printf("%d, %d := %d+direction.X, %d+direction.Y\n", r, c, visitedLocation[1], visitedLocation[0])
		// fmt.Printf("g[%d][%d] = g[%d][%d]\n", c, r, visitedLocation[0], visitedLocation[1])
		g[c][r] = g[visitedLocation[0]][visitedLocation[1]]
		g[visitedLocation[0]][visitedLocation[1]] = EMPTY
	}

	return true
}

func moveDoubleBoxes(d direction, row, col int, grid [][]string) bool {

	if d.value == '>' || d.value == '<' {
		return dfs(d, row, col, grid)
	}
	return bfs(d, row, col, grid)
}

func processInstruction(dir rune, robotRow, robotCol int, g grid) (int, int) {
	direction := directionMap[dir]

	nextRow := robotRow + direction.rowOffset
	nextCol := robotCol + direction.colOffset

	if !g.onGrid(nextRow, nextCol) || g[nextRow][nextCol] == WALL {
		// Robot hasn't moved as move isn't valid move. Either not on grid or it's a wall...
		return robotRow, robotCol
	}

	if g[nextRow][nextCol] == EMPTY {
		// Swap robot with next empty as it's safe
		g[robotRow][robotCol] = EMPTY
		g[nextRow][nextCol] = ROBOT
		return nextRow, nextCol
	}

	if g[nextRow][nextCol] == BOX {
		for {
			// move along, otherwise were checking ourselves
			nextRow += direction.rowOffset
			nextCol += direction.colOffset
			// Loop until we find a wall or empty. Since we are tracking boxes this is basically the end of the current object run
			if g[nextRow][nextCol] == EMPTY {
				// you dont have to move ALL boxes in the chain, only the first
				// If we find an empty we do the following
				// 1. Set the robots current [r][c] to '.'
				// 2. Move the robot 1 unit in current direction, which will override the BOX value at this [r][c]
				// 3. set the empty [r][c] to be a box
				// example
				// @OO.
				// 1 @OO. -> .OO.
				// 2 .OO. -> .@O.
				// 3. .@O. -> .@OO
				// then return the new grid and where the robot is now at
				// This can only work for SINGLE width boxes, it's too simple to 2 widths
				g[robotRow][robotCol] = EMPTY
				g[robotRow+direction.rowOffset][robotCol+direction.colOffset] = ROBOT
				g[nextRow][nextCol] = BOX
				return robotRow + direction.rowOffset, robotCol + direction.colOffset
			}

			if g[nextRow][nextCol] == WALL {
				// Not possible to move, boxes must be up against wall!
				return robotRow, robotCol
			}
		}
	}

	// Part 2 has double width boxes, bit of a PITA so we need to detect if we are at one of those by testing for either side
	if g[nextRow][nextCol] == DOUBLE_BOX_LEFT || g[nextRow][nextCol] == DOUBLE_BOX_RIGHT {
		if !moveDoubleBoxes(direction, nextRow, nextCol, g) {
			// robot didn't move!
			return robotRow, robotCol
		}
		g[nextRow][nextCol] = ROBOT
		g[robotRow][robotCol] = EMPTY
	}

	return nextRow, nextCol
}

func partOne() int {
	grid := createdGrid(gridInput)
	currentRow, currentCol := grid.robotLocation()
	for _, dir := range instructions {
		currentRow, currentCol = processInstruction(dir, currentRow, currentCol, grid)
	}
	return grid.score()
}

func partTwo() int {
	partTwoGridInput := []string{}
	for _, line := range gridInput {
		line = strings.ReplaceAll(line, ".", "..")
		line = strings.ReplaceAll(line, "#", "##")
		line = strings.ReplaceAll(line, "@", "@.")
		line = strings.ReplaceAll(line, "O", "[]")
		partTwoGridInput = append(partTwoGridInput, line)
	}
	grid := createdGrid(partTwoGridInput)
	currentRow, currentCol := grid.robotLocation()
	for _, dir := range instructions {
		currentRow, currentCol = processInstruction(dir, currentRow, currentCol, grid)
	}

	return grid.score()
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())
}
