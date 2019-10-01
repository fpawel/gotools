package main

import (
	"flag"
	"github.com/fpawel/gotools/pkg/rungo"
	"log"
)

func main() {
	log.SetFlags(log.Ltime)

	exeName := flag.String("exe", "", "path to executable")
	args := flag.String("args", "", "command line arguments to pass")
	guiMsgConsole := flag.Int("gui.msg.console", -1, "the code of \"console\" windows message")
	guiMsgPanic := flag.Int("gui.msg.console", -1, "the code of \"panic\" windows message")
	guiWindowClass := flag.String("gui.window.class", "", "window class name to send messages")

	flag.Parse()
	rungo.Cmd{
		ExeName: *exeName,
		ExeArgs: *args,
		NotifyGUI: rungo.NotifyGUI{
			MsgCodeConsole: uintptr(*guiMsgConsole),
			MsgCodePanic:   uintptr(*guiMsgPanic),
			WindowClass:    *guiWindowClass,
		},
	}.Exec()
}
