package main

import (
	"flag"
	"github.com/fpawel/gorunex/pkg/gorunex"
	"log"
)

func main() {
	log.SetFlags(log.Ltime)

	var exeName, args string
	flag.StringVar(&exeName, "exe", "", "path to executable")
	flag.StringVar(&args, "args", "", "command line arguments to pass")

	flag.Parse()

	log.Println("log file:", gorunex.LogFileName())
	if err := gorunex.Process(exeName, args, nil); err != nil {
		log.Fatal(err)
	}
}
