package utils

import (
	"os"
)

func LoadConfiguration(filePath string) error {
	updateConfig := func(path string) error {
		file, err := os.Open(path)
		if err != nil {
			//fmt.Printf("file %v not opened\n", path)
			return err
		} else {
			//fmt.Printf("file %v opened\n", path)
		}
		defer func() {
			closeErr := file.Close()
			if closeErr != nil {
				//fmt.Printf("error on closing file %s\n", path)
			} else {
				//fmt.Printf("file %s successfuly closed\n", path)
			}
		}()
		return nil
	}

	for i := 0; i < 100; i++ {
		err := updateConfig(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
