/*
Tool to export annotation results from the internal DB to an output DB.

Usage: export <path-to-origin.db> <path-to-destination.db>
*/
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

const desc = `Exports annotation results from the internal input database to a new output
database.
The destination file must not exist.`

var opts struct {
	Args struct {
		Input  string `description:"SQLite database filepath"`
		Output string `description:"SQLite database filepath"`
	} `positional-args:"yes" required:"yes"`
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.LongDescription = desc

	if _, err := parser.Parse(); err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}

			fmt.Println()
			parser.WriteHelp(os.Stdout)
		}

		os.Exit(1)
	}

	if _, err := os.Stat(opts.Args.Input); os.IsNotExist(err) {
		log.Fatalf("File %q does not exist", opts.Args.Input)
	}

	// TODO: until we have more information about the output schema needed,
	// this tool does a simple copy&paste of the internal DB. The annotation
	// results are stored in the assignments table.

	if err := copy(opts.Args.Input, opts.Args.Output); err != nil {
		log.Fatal(err)
	}
}

func copy(source, destination string) error {
	src, err := os.Open(source)
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.OpenFile(destination, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	return nil
}
