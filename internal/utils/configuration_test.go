package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	const testFilename = "file"
	err := ioutil.WriteFile(testFilename, []byte(""), 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v", err)
	}

	defer func(){
		err := os.Remove(testFilename);
		if err != nil{
			fmt.Printf("cann't remove test data: %s", testFilename )
		}
	}()
	tests:= []struct{
		name string
		filePath string
		errIsNil bool
	}{
		{
			filePath: testFilename, errIsNil: true,name: "file exists",
		},
		{
			filePath: "file_not_exists", errIsNil: false,name: "file not exists",
		},
	}

	for _, tcase := range tests{
		t.Run(tcase.name, func(t* testing.T){
			assert.Equal(t, LoadConfiguration(tcase.filePath)==nil, tcase.errIsNil)
		})
	}

}
