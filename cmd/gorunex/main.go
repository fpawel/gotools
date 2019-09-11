package main

import (
	"flag"
	"github.com/fpawel/gotools/pkg/run"
	"log"
)

func main() {
	log.SetFlags(log.Ltime)

	var exeName, args string
	flag.StringVar(&exeName, "exe", "", "path to executable")
	flag.StringVar(&args, "args", "", "command line arguments to pass")

	flag.Parse()

	log.Println("log file:", run.LogFileName())
	if err := run.Process(exeName, args, nil); err != nil {
		log.Fatal(err)
	}
}
