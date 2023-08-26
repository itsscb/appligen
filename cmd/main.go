package main

import (
	"log"

	"github.com/itsscb/appligen"
)

func main() {
	a, err := appligen.NewFromFile("data.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = a.Generate("application.tex")
	if err != nil {
		log.Fatal(err)
	}
}
