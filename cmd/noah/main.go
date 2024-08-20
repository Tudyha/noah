package main

import (
	"noah/internal/server"
)

func main() {
	// new server
	s := server.NewServer()

	//start server
	s.Run()
}
