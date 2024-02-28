package main

import (
	"log"

	"github.com/skantay/hezzl/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("app cannot be started: %v", err)
	}
}
