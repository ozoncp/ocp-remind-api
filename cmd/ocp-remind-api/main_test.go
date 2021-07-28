package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfiguration(t *testing.T) {
	const testFilename = "test_file"
	t.Cleanup(func() { _ = os.Remove(testFilename) }) // поставить после константы
	err := ioutil.WriteFile(testFilename, []byte(""), 0600)
	if err != nil {
		t.Fatalf("Unable to create file: %v", err)

	}
	tests := []struct {
		name        string
		filePath    string
		expectedErr bool
	}{
		{
			filePath: testFilename, expectedErr: true, name: "file exists",
		},
		{
			filePath: "file_not_exists", expectedErr: false, name: "file not exists",
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			assert.Equal(t, tCase.expectedErr, LoadConfiguration(tCase.filePath) == nil)
		})
	}

}
