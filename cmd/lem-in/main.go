package main

import (
	"log"

	"lem-in/internal/app"
)

func main() {
	// app.Run - game of the ants ... !)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
