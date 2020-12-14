package main

import (
	"github.com/ExchangeUnion/xud-launcher/cmd"
	"log"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
