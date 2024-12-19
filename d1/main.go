package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

var left []int
var right []int
var rightCopies = make(map[int]int)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), "   ")
		leftId, _ := strconv.Atoi(parts[0])
		left = append(left, leftId)
		rightId, _ := strconv.Atoi(parts[1])
		right = append(right, rightId)
	}
	slices.Sort(left)
	slices.Sort(right)
}

func zip(slices ...[]int) func() []int {
	zippedResult := make([]int, len(slices), len(slices))
	iteration := 0
	return func() []int {
		for ptr := range slices {
			if iteration >= len(slices[ptr]) {
				return nil
			}
			zippedResult[ptr] = slices[ptr][iteration]
		}
		iteration++
		return zippedResult
	}
}

// not using math.abs to avoid 100s of type conversions to float64
func Abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func partOne() int {
	total := 0
	iterator := zip(left, right)
	for pairing := iterator(); pairing != nil; pairing = iterator() {
		total += Abs(pairing[0] - pairing[1])
	}
	return total

}

func partTwo() int {
	total := 0
	var seenCounts = make(map[int]int, len(left))

	for _, leftId := range left {
		count := 0
		if val, ok := seenCounts[leftId]; ok {
			total += val
			continue
		}
		for _, rightId := range right {
			if rightId > leftId {
				break
			}
			if rightId == leftId {
				count++
			}
		}
		leftIdSimilar := leftId * count
		seenCounts[leftId] = leftIdSimilar
		total += leftIdSimilar
	}

	return total
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
