package main

import (
	"log"

	"github.com/shunk031/tcardgen/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
