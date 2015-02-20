package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt.go"
	"github.com/op/go-logging"
)

const configFilename = "./config.toml"

var (
	log = logging.MustGetLogger("nutrition")
)

const usage = `
	Usage:
	nutrition settings show
	nutrition settings set <entry> <value>
`

var config Config

func main() {
	var err error

	args, _ := docopt.Parse(usage, nil, true, "nutrition", false)

	config, err = configRead(configFilename)
	if err != nil {
		log.Fatal(err)
	}

	if args["settings"].(bool) {
		if args["show"].(bool) {
			settingsShow()
			os.Exit(0)
		}

		if args["set"].(bool) {
			settingsSet(
				args["<entry>"].(string),
				args["<value>"].(string),
			)

			os.Exit(0)
		}
	}
}

func settingsShow() {
	fmt.Print(config)
}

func settingsSet(entry string, value string) {
	newConfig, err := configChange(
		config,
		entry,
		value,
	)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = configWrite(configFilename, newConfig)

	if err != nil {
		log.Fatal(err.Error())
	}
}

//@TODO: command line tool interface
//	./cmd add <product>
//	./cmd check <product> <weight>
//	./cmd eat (breakfast|lunch|snack|dinner) <product> <weight>
//  ./cmd journal
//  ./cmd journal (today)
