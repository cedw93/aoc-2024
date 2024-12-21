package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type report struct {
	id     int
	levels []int
}

var reports []report

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	count := 0
	for scanner.Scan() {
		rawReport := scanner.Text()
		rawReportParts := strings.Split(rawReport, " ")
		levels := make([]int, len(rawReportParts), len(rawReportParts))
		for idx, levelId := range rawReportParts {
			levelIdConverted, _ := strconv.Atoi(levelId)
			levels[idx] = levelIdConverted
		}
		reports = append(reports, report{id: count, levels: levels})
		count++
	}

}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func (r report) isSafe() bool {
	decreasing := false
	increasing := false
	tracker := 0
	for idx, level := range r.levels {
		if idx == 0 {
			tracker = level
			continue
		}

		diff := tracker - level
		absDiff := abs(diff)

		if diff == 0 || absDiff > 3 {
			return false
		}

		if diff > 0 {
			if increasing {
				return false
			}
			decreasing = true
		} else {
			if decreasing {
				return false
			}
			increasing = true
		}

		tracker = level

	}
	return true
}

func removeIdx(s []int, idx int) []int {
	copy(s[idx:], s[idx+1:])
	// last element is duplicated with copy so trim it
	return s[:len(s)-1]
}

func partOne() int {
	safeReports := 0
	for _, report := range reports {
		if report.isSafe() {
			safeReports++
		}
	}
	return safeReports
}

func partTwo() int {
	safeReports := 0
	for _, report := range reports {
		if report.isSafe() {
			safeReports++
		} else {
			originalLevels := report.levels
			for idx := range originalLevels {
				if idx != 0 {
					report.levels = originalLevels
				}
				tempLevels := append([]int{}, report.levels...)
				report.levels = removeIdx(tempLevels, idx)
				if report.isSafe() {
					safeReports++
					break
				}
			}
		}
	}
	return safeReports
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
