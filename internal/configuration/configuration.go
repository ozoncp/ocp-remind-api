package configuration

import (
	"errors"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	DB     DB     `yaml:"DB"`
	Grpc   Grpc   `yaml:"Grpc"`
	Kafka  DB     `yaml:"Kafka"`
	Jaeger Jaeger `yaml:"Jaeger"`
}

var c *Configuration

func Instance() *Configuration {
	if c == nil {
		c = &Configuration{}
		return c
	}
	return c
}

func Init(filename string) error {
	if c != nil {
		log.Error().Msg("configuration already inited")
		return errors.New("configuration already inited")
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Err(err).Msg("Cann't open configuration file")
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&c); err != nil {
		log.Err(err).Msg("Cann't decode configuration file")
		return err
	}

	return nil
}

type DB struct {
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbName"`
	SslMode  string `yaml:"sslmode"`
}

func (db DB) URI() string {
	//postres://postgres:postgres@localhost:5432/reminds?sslmode=none
	return "postgres://" +
		db.UserName + ":" + db.Password + "@" +
		db.Host + ":" + db.Port + "/" + db.DBName + "?sslmode=none"
}

type Grpc struct {
	Port string `yaml:"port"`
}

type Kafka struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func (k Kafka) URI() string {
	return k.Host + ":" + k.Port
}

type Jaeger struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
