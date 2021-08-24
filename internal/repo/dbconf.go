package repo

import (
	"encoding/json"
	"fmt"
	"os"
)

type DbConfiguration struct {
	UserName string
	Passwd   string
	Host     string
	Port     string
	DbName   string
	Options  []string
}

func (dc DbConfiguration) DSN() string {
	//"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	dsn := fmt.Sprintf("postgres://%s:%s@:%s:%s/%s",
		dc.UserName, dc.Passwd, dc.Host, dc.Port, dc.DbName)
	if len(dc.Options) != 0 {
		dsn += "?"
		for i := 0; i+1 < len(dc.Options); i++ {
			dsn += dc.Options[i] + ","
		}
		dsn += dc.Options[len(dc.Options)-1]
	}
	return dsn
}

func (dc *DbConfiguration) ReadConfiguration(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(dc)
	if err != nil {
		return err
	}
	return nil
}
