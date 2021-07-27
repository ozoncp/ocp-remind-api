package main

import (
	"fmt"
	"os"
)

func LoadConfiguration(filePath string) error {
	for i := 0; i < 100; i++ {
		err := func() error {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			return nil
		}()
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	fmt.Println("ocp remind api project")
}
