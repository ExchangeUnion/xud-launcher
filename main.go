package main

import (
	"github.com/reliveyy/xud-launcher/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
