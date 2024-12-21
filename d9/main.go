package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type block struct {
	blockId    int
	startIndex int
	endIndex   int
	size       int
	free       bool
}

var files = make(map[int]*block)
var blanks = []*block{}
var largestFileId = -1
var totalLineSize = 0

func init() {
	scanner := bufio.NewScanner(os.Stdin)
	fileId := 0
	processedBlocks := 0
	currentIndex := 0
	for scanner.Scan() {
		data := scanner.Text()
		for _, r := range data {
			runeAsInt := int(r - '0')
			totalLineSize += runeAsInt
			if processedBlocks%2 == 0 {
				// even means were on a file block...
				files[fileId] = &block{fileId, currentIndex, currentIndex + runeAsInt - 1, runeAsInt, false}
				largestFileId = fileId
				fileId++
			} else {
				// odds are blank
				blanks = append(blanks, &block{-1, currentIndex, currentIndex + runeAsInt - 1, runeAsInt, true})
			}
			currentIndex += runeAsInt
			processedBlocks++
		}
	}
}

func (b *block) checkSum() int {
	result := 0
	for i := 0; i < b.size; i++ {
		result += b.blockId * (b.startIndex + i)
	}
	return result
}

func getNextSwappable(currentSwappable int, fileSystem []int) int {
	for i := currentSwappable - 1; i >= 0; i-- {
		if fileSystem[i] != -1 {
			return i
		}
	}
	return -1
}

func checkSum(fileSystem []int) int {
	result := 0
	for i, v := range fileSystem {
		if v != -1 {
			result += i * v
		}
	}
	return result
}

func partOne() int {
	// bit hacky but saves messing up part 2s data! could be refactored for sure, horrendous time complexity for sure
	// fmt.Println("Line size", totalLineSize)

	result := make([]int, totalLineSize, totalLineSize)
	blankIndexes := []int{}

	for _, blank := range blanks {
		for i := 0; i < blank.size; i++ {
			result[blank.startIndex+i] = -1
			blankIndexes = append(blankIndexes, blank.startIndex+i)
		}
	}

	for fileId, block := range files {
		for i := 0; i < block.size; i++ {
			result[block.startIndex+i] = fileId
		}
	}

	nextSwappable := getNextSwappable(len(result), result)

	for _, blankIndex := range blankIndexes {
		if blankIndex > nextSwappable {
			break
		}
		result[blankIndex], result[nextSwappable] = result[nextSwappable], result[blankIndex]
		nextSwappable = getNextSwappable(nextSwappable, result)
	}

	return checkSum(result)
}

func partTwo() int {
	result := 0
	currentFileId := largestFileId

	for currentFileId > -1 {
		fileBlock := files[currentFileId]
		for _, blank := range blanks {
			if !blank.free {
				continue
			}
			if blank.startIndex >= fileBlock.startIndex {
				// We've gone too far as files can only move left
				break
			}
			if fileBlock.size <= blank.size {
				// fmt.Printf("Can Move file id %d into blank %+v\n", currentFileId, blank)
				// Now we can move file...
				fileBlock.startIndex = blank.startIndex
				fileBlock.endIndex = blank.startIndex + fileBlock.size - 1
				blank.size = blank.size - fileBlock.size
				// Not -1 here because we want the next free start place not the end of the old blank
				// e.g consider ..1
				// if blank is current;y
				// start 0 end 1 size 2
				// and we move 1 into it (with size 1)
				// the blank is now
				// start 0 + 1, size 2-1 (blank.size - file.size)
				// Doing minus -1 before incrementing index would mean start would be 0, which is where we just inserted
				// essentially this means we set startIndex to lastInserted +!
				blank.startIndex = blank.startIndex + fileBlock.size
				if blank.size == 0 {
					blank.free = false
				}
				break
			}
		}
		currentFileId--
	}

	for _, file := range files {
		result += file.checkSum()
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
