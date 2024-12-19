package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type directionMap struct {
	rowOffset int
	colOffset int
	label     string
}

type crossConfig struct {
	first  []directionMap
	second []directionMap
}

var wordGrid [][]rune
var xmas = [...]rune{'X', 'M', 'A', 'S'}
var directions = [...]directionMap{
	{0, 1, "right"},
	{1, 0, "down"},
	{0, -1, "left"},
	{-1, 0, "up"},
	{-1, -1, "upLeft"},
	{-1, 1, "upRight"},
	{1, -1, "downLeft"},
	{1, 1, "downRight"},
}

var crossMappings = [...][]directionMap{
	{
		{-1, -1, "upLeft"},
		{1, 1, "downRight"},
	},
	{
		{1, -1, "downLeft"},
		{-1, 1, "upRight"},
	},
}

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		wordGrid = append(wordGrid, []rune(line))
	}
}

func withinGrindBoundary(row, col int) bool {
	return row > -1 && col > -1 && row < len(wordGrid) && col < len(wordGrid[row])
}

func calcWordsFromX(row, col int) int {
	wordsFound := 0
	for _, direction := range directions {
		for i := 1; i < len(xmas); i++ {
			offsetRow := row + direction.rowOffset*i
			offsetColumn := col + direction.colOffset*i
			if !withinGrindBoundary(offsetRow, offsetColumn) {
				break
			}

			if wordGrid[offsetRow][offsetColumn] != xmas[i] {
				break
			}

			// prevents partial matches on line ends triggering being counted
			if i == 3 {
				wordsFound += 1
			}
		}
	}
	return wordsFound
}

func isXPattern(row, col int) bool {
	for _, mapping := range crossMappings {
		pairSum := 0
		for _, direction := range mapping {
			offsetRow := row + direction.rowOffset
			offsetColumn := col + direction.colOffset
			if !withinGrindBoundary(offsetRow, offsetColumn) {
				return false
			}
			pairSum += int(wordGrid[offsetRow][offsetColumn])
		}

		// ensuring they are in the correct pairing is handled using the mappings already
		// We just care than any pairing is M and S since we only check this when we are on an A
		// If this pair doesn't match then no point checking the other
		if pairSum != int(xmas[1])+int(xmas[3]) {
			return false
		}
	}
	return true
}

func partOne() int {
	result := 0
	for rowId, row := range wordGrid {
		for colId, c := range row {
			if c == xmas[0] {
				result += calcWordsFromX(rowId, colId)
			}
		}
	}

	return result
}

func partTwo() int {
	result := 0
	for rowId, row := range wordGrid {
		for colId, c := range row {
			// We care about A as its the middle of the cross
			if c == xmas[2] {
				if isXPattern(rowId, colId) {
					result += 1
				}
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
