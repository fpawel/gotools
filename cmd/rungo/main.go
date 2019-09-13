package main

import (
	"flag"
	"github.com/fpawel/gotools/pkg/rungo"
	"log"
)

func main() {
	log.SetFlags(log.Ltime)

	exeName := flag.String( "exe", "", "path to executable")
	args := flag.String( "args", "", "command line arguments to pass")
	useGui := flag.Bool( "gui", true, "use GUI (true|false)")

	flag.Parse()
	rungo.Cmd{
		ExeName:   *exeName,
		ExeArgs:   *args,
		UseGUI:    false,
		NotifyGUI: rungo.NotifyGUI{
			MsgCodeConsole: 0,
			MsgCodePanic:   0,
			WindowClass:    "",
		},
	}.Exec()

	if .; err != nil {
		log.Fatal(err)
	}
}
