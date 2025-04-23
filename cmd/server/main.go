package main

import (
	"github.com/gekich/go-web/internal/server"
)

func main() {
	s := server.New()
	s.Init()
	s.Run()
}
