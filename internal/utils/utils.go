package utils

var filter = map[int]struct{}{
	3:  {},
	5:  {},
	7:  {},
	11: {},
	13: {},
	17: {},
	19: {},
}

func Filter(input []int) []int {
	for i := 0; i < len(input); i++ {
		if _, ok := filter[input[i]]; ok {
			input[i] = input[len(input)-1]
			input = input[:len(input)-1]
			i -= 1
		}
	}
	return input
}

func Batch(input []int, size int) [][]int {
	result := make([][]int, (len(input)+size-1)/size)
	pos, counter := 0, 0
	for pos < len(input) {
		if len(input)-pos < size {
			size = len(input) - pos
		}
		result[counter] = input[pos : pos+size]
		pos += size
		counter++
	}
	return result
}

func Mirror(source map[string]int) map[int]string {
	if len(source) == 0 {
		return map[int]string{}
	}
	result := make(map[int]string, len(source))
	for key, value := range source {
		result[value] = key
	}
	return result
}
