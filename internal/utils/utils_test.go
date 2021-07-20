package utils

import "testing"

func TestMirrorMapSimple(t *testing.T) {
	sourceMap := make(map[string]int)
	sourceMap["one"] = 1
	sourceMap["two"] = 2
	sourceMap["three"] = 3
	convertedMap := MirrorMap(sourceMap)
	if len(convertedMap) != len(sourceMap) {
		t.Error("dest map size differ")
	}
	if convertedMap[1] != "one" {
		t.Error("Error")
	}
}

func TestGetBatchesCount(t *testing.T) {
	if getChunksCount(0, 5) != 0 {
		t.Error("Wrong batches count for empty zero slice size")
	}
	if getChunksCount(5, 5) != 1 {
		t.Error("Wrong batches size [5, 5]")
	}
	if getChunksCount(5, 4) != 2 {
		t.Error("Wrong batches size [5, 4]")
	}
	if getChunksCount(100, 1) != 100 {
		t.Error("Wrong batches size [100, 1]")
	}
}

func TestMin(t *testing.T) {
	if min(-1, 1) != -1 {
		t.Error("[-1, 1]")
	}
	if min(1, 1) != 1 {
		t.Error("[1, 1]")
	}
	if min(-4, -7) != -7 {
		t.Error("[-4, -7]")
	}
	if min(4, 7) != 4 {
		t.Error("[4, 7]")
	}

}

func TestCutSlice(t *testing.T) {
	source := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	res := BatchSlice(source, 3)
	if len(res) != 4 {
		t.Error("wrong output slice size")
	}
	testRes := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}}
	for i := 0; i < len(res); i++ {
		for j := 0; j < len(res[i]); j++ {
			if res[i][j] != testRes[i][j] {
				t.Error("Wrong result")
			}
		}
	}
}

func TestFilterSlice(t *testing.T) {
	{
		filter := map[int]struct{}{3: {}, 5: {}, 7: {}, 11: {}, 13: {}, 17: {}, 19: {}}
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
		result := FilterSlice(slice, filter)
		if len(result) != 13 {
			t.Error("Wrong filtered slice size")
		}
		testSlice := []int{1, 2, 4, 6, 8, 9, 10, 12, 14, 15, 16, 18, 20}
		for i := 0; i < len(result); i++ {
			if testSlice[i] != result[i] {
				t.Fatal("Wrong filtered slice contents")
			}
		}
	}
	{
		filter := map[int]struct{}{}
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
		result := FilterSlice(slice, filter)
		if len(result) != len(slice) {
			t.Error("Wrong result for empty filter")
		}
	}
	{
		filter := map[int]struct{}{3: {}, 5: {}, 7: {}, 11: {}, 13: {}, 17: {}, 19: {}}
		slice := []int{}
		result := FilterSlice(slice, filter)
		if len(result) != 0 {
			t.Error("Wrong result for empty slice")
		}
	}

}
