package utils

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

func FilterSlice(inputSlice []int, filter map[int]struct{}) []int {
	result := make([]int, 0, len(inputSlice))
	for i := range inputSlice {
		if _, ok := filter[inputSlice[i]]; !ok {
			result = append(result, inputSlice[i])
		}
	}
	return result
}

func BatchSlice(slice []int, chunkSize int) [][]int {
	chunksCount := getChunksCount(len(slice), chunkSize)
	resultSlice := make([][]int, chunksCount)
	pos, chunkCounter, copySize := 0, 0, 0
	for pos < len(slice) {
		copySize = min(len(slice)-pos, chunkSize)
		resultSlice[chunkCounter] = make([]int, copySize)
		copy(resultSlice[chunkCounter], slice[pos:pos+copySize])
		pos += copySize
		chunkCounter++
	}
	return resultSlice
}

func MirrorMap(sourceMap map[string]int) map[int]string {
	outMap := make(map[int]string, len(sourceMap))
	for key, value := range sourceMap {
		outMap[value] = key
	}
	return outMap
}
