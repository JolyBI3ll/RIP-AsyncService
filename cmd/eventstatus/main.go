package main

import (
	"log"
	"async/internal/api"
)

func main() {
	log.Println("App start")
	api.StartServer()
	log.Println("App stop")
}