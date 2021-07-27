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
	var n int
	for _, v := range input {
		if _, ok := filter[v]; !ok {
			input[n] = v
			n++
		}
	}

	return input[:n]
}

func Batch(input []int, size int) [][]int {
	if len(input) <= size {
		return [][]int{input}
	}

	result := make([][]int, 0, (len(input)+size-1)/size)
	for i, j := 0, size; j <= len(input); i, j = i+size, j+size {
		result = append(result, input[i:j])
	}
	if v := len(input) % size; v != 0 {
		result = append(result, input[len(input)-v:])
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
