package main

import (
	"log"

	"github.com/rm-ryou/mococoplan/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
