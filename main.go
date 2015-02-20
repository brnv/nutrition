package main

import (
	"fmt"

	"github.com/docopt/docopt.go"
	"github.com/op/go-logging"
)

const configFilename = "./config.toml"

var (
	log    = logging.MustGetLogger("nutrition")
	mode   = ""
	action = ""
)

const usage = `
	Usage:
	nutrition settings show
`

func main() {
	args, _ := docopt.Parse(usage, nil, true, "nutrition", false)

	config, _ := configRead(configFilename)

	if _, ok := args["settings"]; ok {
		mode = "settings"

		if _, ok := args["show"]; ok {
			action = "show"
		}
	}

	switch mode {
	case "settings":
		if action == "show" {
			fmt.Print(config)
		}

	}

}

//@TODO: update configuration file with new data
//	./cmd settings set <entry> <value>

//@TODO: command line tool interface
//	./cmd add <product>
//	./cmd check <product> <weight>
//	./cmd eat (breakfast|lunch|snack|dinner) <product> <weight>
//  ./cmd journal
//  ./cmd journal (today)
