package main

import (
	"fmt"
	"os"
)

func LoadConfiguration(filePath string) error {
	openFile := func() error {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer func() {
			err := file.Close()
			if err != nil {
				fmt.Println("error on file closing...")
			}
		}()
		return nil
	}
	for i := 0; i < 10; i++ {
		err := openFile()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fmt.Println("ocp remind api project")
}
