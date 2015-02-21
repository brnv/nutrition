package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

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
