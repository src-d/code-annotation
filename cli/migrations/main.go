package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/jessevdk/go-flags"
)

const (
	description string = "Migrate internal database"
	sqliteDSN          = "sqlite://%s"
)

func main() {
	parser := flags.NewParser(nil, flags.Default)
	addCommand(parser, &uastColsOpts)
	addCommand(parser, &uastImportOpts)
	addCommand(parser, &diffRmColOpts)
	addCommand(parser, &featureDropOpts)
	addCommand(parser, &vacuumOpts)
	parse(parser, description)
}

func addCommand(parser *flags.Parser, command Command) {
	if _, err := parser.AddCommand(command.Name(), command.ShortDesc(), command.LongDesc(), command); err != nil {
		panic(err)
	}
}

func parse(parser *flags.Parser, description string) {
	parser.LongDescription = description
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
}
