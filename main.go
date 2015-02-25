package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"text/template"

	"github.com/docopt/docopt.go"
	"github.com/op/go-logging"
	"github.com/seletskiy/tplutil"
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

const productImpactTpl = `
	carbohydrates = {{.Carbohydrates}}%{{"\n"}}
	proteins = {{.Proteins}}%{{"\n"}}
	fats = {{.Fats}}%{{"\n"}}
	calories = {{.Calories}}%{{"\n"}}
`

func checkProductImpact(productName string, weight float64) {
	products, _ := productsRead(productsFile)

	for _, product := range products.Product {
		if product.Name == productName {
			impact := struct {
				Carbohydrates float64
				Proteins      float64
				Fats          float64
				Calories      float64
			}{
				Carbohydrates: product.Carbohydrates * weight / 100 / config.Settings.Carbohydrates * 100,
				Proteins:      product.Proteins * weight / 100 / config.Settings.Proteins * 100,
				Fats:          product.Fats * weight / 100 / config.Settings.Fats * 100,
				Calories:      product.Calories * weight / 100 / config.Settings.Calories * 100,
			}

			myTpl := template.Must(
				template.New("productImpactTpl").Parse(tplutil.Strip(
					productImpactTpl,
				)))

			buf := bytes.NewBuffer([]byte{})
			myTpl.Execute(buf, impact)

			fmt.Print(buf.String())
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
//	./cmd product list
//	./cmd product show <product>
//	./cmd product edit
//	./cmd product edit <product>
//  ./cmd journal (edit) (today)
