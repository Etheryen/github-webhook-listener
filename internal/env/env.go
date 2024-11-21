package env

import (
	"os"

	"github.com/joho/godotenv"
)

var requiredVars = []string{"GITHUB_SECRET", "BRANCH", "COMMAND", "PORT"}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("failed to load .env")
	}

	for _, envVar := range requiredVars {
		if os.Getenv(envVar) == "" {
			panic(envVar + " not set in .env")
		}
	}
}
