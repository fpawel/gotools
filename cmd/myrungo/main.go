package main

import (
	"flag"
	"github.com/fpawel/gotools/pkg/rungo"
	"log"
)

func main() {
	log.SetFlags(log.Ltime)

	var exeName, args string
	flag.StringVar(&exeName, "exe", "", "path to executable")
	flag.StringVar(&args, "args", "", "command line arguments to pass")

	flag.Parse()

	log.Println("log file:", rungo.LogFileName())
	if err := rungo.Process(exeName, args, nil); err != nil {
		log.Fatal(err)
	}
}
