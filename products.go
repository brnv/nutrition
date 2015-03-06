package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/seletskiy/tplutil"
)

//@TODO: product nutition fact to separate structure
type product struct {
	Name          string
	Carbohydrates float64 `toml:"carbohydrates"`
	Proteins      float64 `toml:"proteins"`
	Fats          float64 `toml:"fats"`
	Calories      float64 `toml:"calories"`
}

type Products struct {
	Product []product `toml:"product"`
}

const productImpactTpl = `
	c {{printf "%.2f" .Carbohydrates}}%, p {{printf "%.2f" .Proteins}}%,
	f {{printf "%.2f" .Fats}}%, cal = {{printf "%.2f" .Calories}}%
`

func productAdd(products Products, productName string) (Products, error) {
	newProduct := Products{
		Product: []product{
			product{
				Name: productName,
			},
		},
	}

	tmpFilename := "/tmp/nutrition-add-product-" + productName + ".toml"

	err := productsWrite(tmpFilename, newProduct)
	if err != nil {
		return products, err
	}

	cmd := exec.Command("/bin/vim", tmpFilename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return products, err
	}

	newProduct, err = productsRead(tmpFilename)
	if err != nil {
		return products, err
	}

	//@TODO: [0] looks ugly, find another way to do that
	products.Product = append(products.Product, newProduct.Product[0])

	return products, nil
}

func productsRead(filename string) (Products, error) {
	products := Products{}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return Products{}, err
	}

	if _, err := toml.Decode(string(contents), &products); err != nil {
		return Products{}, err
	}

	return products, nil
}

func productsWrite(filename string, products Products) error {
	buf := bytes.NewBuffer([]byte{})
	tomlEncoder := toml.NewEncoder(buf)
	tomlEncoder.Encode(products)

	err := ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func productImpact(productName string, weight float64) string {
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

			return buf.String()
		}
	}

	return ""
}
