package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ClientEnvs struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func LoadClientEnvs() ClientEnvs {
	var errors []string

	err := godotenv.Load("./.client.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host, ok := os.LookupEnv("DB_HOST")
	if !ok {
		errors = append(errors, "DB_HOST is not set")
	}

	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		errors = append(errors, "DB_PORT is not set")
	}

	username, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		errors = append(errors, "DB_USERNAME is not set")
	}

	password, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		errors = append(errors, "DB_PASSWORD is not set")
	}

	database, ok := os.LookupEnv("DB_DATABASE")
	if !ok {
		errors = append(errors, "DB_DATABASE is not set")
	}

	if len(errors) > 0 {
		for _, e := range errors {
			log.Println(e)
		}

		os.Exit(1)
	}

	return ClientEnvs{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}

}
