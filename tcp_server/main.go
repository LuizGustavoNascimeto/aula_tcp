package main

import (
	"log"
	app "tcp_server/src/server"
)

func main() {
	server := app.NewServer(":8080")
	if err := server.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
