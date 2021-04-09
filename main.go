package main

import (
	"log"

	"github.com/arata-nvm/offblast/config"
	"github.com/arata-nvm/offblast/web"
)

func main() {
	s := web.NewServer()

	if err := s.Start(":" + config.Port()); err != nil {
		log.Fatalln(err)
	}
}
