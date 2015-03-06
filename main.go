package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/docopt/docopt.go"
	"github.com/op/go-logging"
)

//@TODO: rename to configFile
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
	nutrition check <product> <weight>
	nutrition eat (breakfast|lunch|snack|dinner) <product> <weight>
	nutrition journal (list|today)
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

	if args["check"].(bool) {
		weightFloat, err := strconv.ParseFloat(args["<weight>"].(string), 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		checkProductImpact(
			args["<product>"].(string),
			weightFloat,
		)
	}

	if args["eat"].(bool) {
		mealType := ""
		if args["breakfast"].(bool) {
			mealType = "breakfast"
		}

		if args["lunch"].(bool) {
			mealType = "lunch"
		}

		if args["snack"].(bool) {
			mealType = "snack"
		}

		if args["dinner"].(bool) {
			mealType = "dinner"
		}

		weightFloat, err := strconv.ParseFloat(args["<weight>"].(string), 64)
		if err != nil {
			log.Fatal(err.Error())
		}

		eat(
			mealType,
			args["<product>"].(string),
			weightFloat,
		)
	}

	if args["journal"].(bool) {
		if args["today"].(bool) {
			journalShow("today")
		}

		if args["list"].(bool) {
			journalShow("list")
		}
	}
}

//@TODO: check if this product exists
func eat(mealType string, productName string, weight float64) {
	err := journalAdd(
		mealType,
		productName,
		weight,
	)

	if err != nil {
		log.Fatal(err.Error())
	}
}

//@TODO: move calculations to another file

func checkProductImpact(productName string, weight float64) {
	fmt.Print(productImpact(productName, weight))
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
//	./cmd product list
//	./cmd product show <product>
//	./cmd product edit
//	./cmd product edit <product>
//  ./cmd journal (edit) (today)
