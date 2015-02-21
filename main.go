package main

import (
	"fmt"
	"os"

	"github.com/docopt/docopt.go"
	"github.com/op/go-logging"
)

const configFilename = "./config.toml"
const productsFile = "./products.toml"

var (
	log = logging.MustGetLogger("nutrition")
)

const usage = `
	Usage:
	nutrition settings show
	nutrition settings set <entry> <value>
	nutrition product add <product>
`

var config Config

func main() {
	var err error

	args, _ := docopt.Parse(usage, nil, true, "nutrition", false)

	config, err = configRead(configFilename)
	if err != nil {
		log.Fatal(err.Error())
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

	//@TODO: check if product with that name exists
	if args["product"].(bool) {
		if args["add"].(bool) {
			products, err := productsRead(productsFile)

			if err != nil {
				log.Fatal(err.Error())
			}

			newProducts, err := productAdd(products, args["<product>"].(string))

			if err != nil {
				log.Fatal(err.Error())
			}

			err = productsWrite(productsFile, newProducts)

			if err != nil {
				log.Fatal(err.Error())
			}

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
//	./cmd product edit
//	./cmd product edit <product>
//	./cmd check <product> <weight>
//	./cmd eat (breakfast|lunch|snack|dinner) <product> <weight>
//  ./cmd journal
//  ./cmd journal (today)
