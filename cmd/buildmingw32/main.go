package main

import (
	"github.com/powerman/must"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	must.AbortIf = must.PanicIf
	log.SetPrefix("buildmingw32: ")
	log.SetFlags(0)
	must.AbortIf(os.Setenv("GOARCH", "386"))
	must.AbortIf(os.Setenv("CGO_ENABLED", "1"))
	setMinGW32Path()
	printArgs()
	cmd := exec.Command(os.Args[1], os.Args[2:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatalln(err)
	}

	log.Println("ok")
}

func setMinGW32Path() {
	xs := strings.Split(os.Getenv("PATH"), ";")
	sort.Slice(xs, func(i, j int) bool {
		return strings.Compare(xs[i], xs[j]) < 0
	})

	for i, s := range xs {
		if filepath.Base(s) != "bin" {
			continue
		}
		if strings.ToLower(filepath.Base(filepath.Dir(s))) != "mingw" {
			continue
		}
		xs[i] = path.Join(
			filepath.Dir(filepath.Dir(s)), "MinGW32", "bin")
		must.AbortIf(os.Setenv("PATH", strings.Join(xs, ";")))
		log.Println("mingw -> MinGW32")
		return
	}
	log.Fatalln("mingw", ": not found in path:", os.Getenv("PATH"))
}

func printArgs() {
	args := make([]interface{}, len(os.Args[1:]))
	for i, v := range os.Args[1:] {
		args[i] = v
	}
	log.Println(args...)
}
