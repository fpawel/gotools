package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func main() {

	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("sqlstr:", pathS)

	r := regexp.MustCompile(`db_([\w\d_]+)\.sql`)

	_ = filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if f == nil || f.IsDir() {
			return nil
		}
		xs := r.FindStringSubmatch(f.Name())
		if len(xs) == 0 {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		s := xs[1]

		fileName := "sql_" + s + "_generated.go"

		fmt.Println("+", f.Name(), ":")
		fmt.Println("\t", fileName)

		fs, err := os.Create(filepath.Join(filepath.Dir(path), fileName))
		if err != nil {
			panic(err)
		}
		defer fs.Close()

		_, _ = fmt.Fprintln(fs, "package", filepath.Base(filepath.Dir(path)))
		_, _ = fmt.Fprintln(fs, "")
		_, _ = fmt.Fprintf(fs, "const SQL%s = `\n", strcase.ToCamel(s))
		_, _ = fs.Write(b)
		_, _ = fmt.Fprintln(fs, "`")

		return nil
	})

}
