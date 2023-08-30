package main

import (
	"github.com/joho/godotenv"
	"user-segmentation-api/utils/http"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}

	server := http.NewServer(":8000")
	server.StartServer()

}
