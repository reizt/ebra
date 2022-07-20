package conf

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	workdir := os.Getenv("WORKDIR")
	err := godotenv.Load(workdir + "/.env")
	if err != nil {
		panic(err)
	}
}
