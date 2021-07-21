package utils

import "errors"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getChunksCount(sliceSize, batchSlice int) int {
	result := sliceSize / batchSlice
	if sliceSize%batchSlice != 0 {
		result++
	}
	return result
}

var filter = map[int]struct{}{
	3:  {},
	5:  {},
	7:  {},
	11: {},
	13: {},
	17: {},
	19: {},
}

func FilterSlice(inputSlice []int) []int {
	for i := 0; i < len(inputSlice); i++ {
		if _, ok := filter[inputSlice[i]]; ok {
			inputSlice[i] = inputSlice[len(inputSlice)-1]
			inputSlice = inputSlice[:len(inputSlice)-1]
			i -= 1
		}
	}
	return inputSlice
}

func BatchSlice(slice []int, chunkSize int) [][]int {
	chunksCount := getChunksCount(len(slice), chunkSize)
	resultSlice := make([][]int, chunksCount)
	pos, chunkCounter, copySize := 0, 0, 0
	for pos < len(slice) {
		copySize = min(len(slice)-pos, chunkSize)
		resultSlice[chunkCounter] = slice[pos : pos+copySize]
		pos += copySize
		chunkCounter++
	}
	return resultSlice
}

func MirrorMap(sourceMap map[string]int) (map[int]string, error) {
	if sourceMap == nil {
		return nil, errors.New("source map is nil")
	}
	outMap := make(map[int]string, len(sourceMap))
	for key, value := range sourceMap {
		outMap[value] = key
	}
	return outMap, nil
}
