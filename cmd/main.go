package main

import (
	"flag"
	"fmt"
	"github.com/fpawel/gorunex/pkg/gorunex"
	"github.com/maruel/panicparse/stack"
	"io"
	"log"
	"os/exec"
	"strings"
)


func main(){
	log.SetFlags(log.Ltime)

	var exeName, args string
	flag.StringVar(&exeName, "exe", "", "path to executable")
	flag.StringVar(&args, "args", "", "command line arguments to pass")

	flag.Parse()

	log.Println("log file:", gorunex.LogFileName())

	out := NewOutput()

	defer func() {
		log.Println("close log file: ", gorunex.LogFileName(), out.Close())
	}()

	cmd := exec.Command(exeName, strings.Fields(args)...)
	cmd.Stderr = out
	cmd.Stdout = out
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	err := cmd.Wait()
	if err == nil {
		return
	}
	if _,err := fmt.Fprintln(out, err); err != nil {
		log.Fatal(err)
	}

	if err := out.PrintPanic(); err != nil {
		log.Fatal(err)
	}
}


func parseDump(in io.Reader, out io.Writer) error {
	// Optional: Check for GOTRACEBACK being set, in particular if there is only
	// one goroutine returned.
	c, err := stack.ParseDump(in, out, true)
	if err != nil {
		return err
	}

	// Find out similar goroutine traces and group them into buckets.
	buckets := stack.Aggregate(c.Goroutines, stack.AnyValue)

	// Calculate alignment.
	srcLen := 0
	pkgLen := 0
	for _, bucket := range buckets {
		for _, line := range bucket.Signature.Stack.Calls {
			if l := len(line.SrcLine()); l > srcLen {
				srcLen = l
			}
			if l := len(line.Func.PkgName()); l > pkgLen {
				pkgLen = l
			}
		}
	}

	for _, bucket := range buckets {
		// Print the goroutine header.
		extra := ""
		if s := bucket.SleepString(); s != "" {
			extra += " [" + s + "]"
		}
		if bucket.Locked {
			extra += " [locked]"
		}
		if c := bucket.CreatedByString(false); c != "" {
			extra += " [Created by " + c + "]"
		}
		if _, err := fmt.Fprintf(out, "%d: %s%s\n", len(bucket.IDs), bucket.State, extra); err != nil {
			return err
		}

		// Print the stack lines.
		for _, line := range bucket.Stack.Calls {
			if _, err := fmt.Fprintf(out,
				"    %-*s %-*s %s(%s)\n",
				pkgLen, line.Func.PkgName(), srcLen, line.SrcLine(),
				line.Func.Name(), &line.Args); err != nil {
				return err
			}
		}
		if bucket.Stack.Elided {
			if _, err := fmt.Fprintf(out, "    (...)\n"); err != nil {
				return err
			}
		}
	}
	return nil
}



