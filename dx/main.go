package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
	}
}

func timer() func() {
	start := time.Now()
	return func() {
		fmt.Printf("took %v\n", time.Since(start))
	}
}

func partOne() int {
	return 0
}

func partTwo() int {
	return 0
}

func main() {
	defer timer()()
	fmt.Println("Part One:", partOne())
	fmt.Println("Part Two:", partTwo())

}
