package main

import (
	"github.com/joho/godotenv"
	"os"
	"user-segmentation-api/utils/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	server := http.New(os.Getenv("SERVER_ADDRESS"))
	server.StartServer()

}
