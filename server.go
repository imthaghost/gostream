package main

import "github.com/imthaghost/gostream/server"

func main() {

	s := server.NewServer()
	s.Start(":8000")
}
