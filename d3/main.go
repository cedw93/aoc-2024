package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	opPattern          = `mul\((\d+),(\d+)\)`
	instructionPattern = `(?:^|do)[^d]*`
)

var lines []string
var opRe *regexp.Regexp
var instructionRe *regexp.Regexp

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	opRe = regexp.MustCompile(opPattern)
	instructionRe = regexp.MustCompile(instructionPattern)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
}

func aToIIgnoreError(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}

func partOne() int {
	result := 0
	// instructionPattern := `do([^d]*)`
	// instructionRe := regexp.MustCompile(instructionPattern)
	for _, line := range lines {
		matches := opRe.FindAllStringSubmatch(line, -1)
		for _, opMatch := range matches {
			result += aToIIgnoreError(opMatch[1]) * aToIIgnoreError(opMatch[2])
		}
	}
	return result
}

func partTwo() int {
	result := 0
	ignore := false
	for _, line := range lines {
		instructionMatches := instructionRe.FindAllStringSubmatch(line, -1)
		// iMatch[0] will be between 'do's or (or ^ or $)
		// example: xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))
		// would have 3 matches
		// 	1: xmul(2,4)&mul[3,7]!^
		// 	2 don't()_mul(5,5)+mul(32,64](mul(11,8)un
		// 	3 do()?mul(8,5))
		for _, iMatch := range instructionMatches {
			if strings.HasPrefix(iMatch[0], "don't") {
				ignore = true
				continue
			}
			if strings.HasPrefix(iMatch[0], "do") {
				ignore = false
			}

			if !ignore {
				// opMatch 0 will be every occurance of multi(x,y) within match[0]
				// x will be opMatch[1] y will be opMatch[2]
				for _, opMatch := range opRe.FindAllStringSubmatch(iMatch[0], -1) {
					result += aToIIgnoreError(opMatch[1]) * aToIIgnoreError(opMatch[2])
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
