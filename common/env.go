package common

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	prod := os.Getenv("GO_ENV")

	if prod == "development" {
		err := godotenv.Load()
		if err != nil {
			return err
		}
	}

	return nil
}
