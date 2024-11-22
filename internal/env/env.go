package env

import (
	"os"

	"github.com/joho/godotenv"
)

var requiredVars = []string{"GITHUB_SECRET", "PORT", "CONFIG_PATH"}

func GetVars() (string, string, string) {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env")
	}

	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			panic(envVar + " not set in .env")
		}
	}

	return os.Getenv(
			"GITHUB_SECRET",
		), os.Getenv(
			"PORT",
		), os.Getenv(
			"CONFIG_PATH",
		)
}
