package main

import (
	"log"

	"github.com/skantay/service-2/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("app cannot be started: %v", err)
	}
}
