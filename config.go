package main

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/seletskiy/tplutil"
)

type Config struct {
	Settings struct {
		Carbohydrates float64 `json:"carbohydrates"`
		Proteins      float64 `json:"proteins"`
		Fats          float64 `json:"fats"`
		Calories      float64 `json:"calories"`
	}
}

var config Config

const settingsShowTpl = `
	carbohydrates = {{.Settings.Carbohydrates}}{{"\n"}}
	proteins = {{.Settings.Proteins}}{{"\n"}}
	fats = {{.Settings.Fats}}{{"\n"}}
	calories = {{.Settings.Calories}}{{"\n"}}
`

func (config Config) String() string {
	myTpl := template.Must(template.New("show-settings").Parse(tplutil.Strip(
		settingsShowTpl,
	)))

	buf := bytes.NewBuffer([]byte{})
	myTpl.Execute(buf, config)

	return buf.String()
}

func configRead(filename string) (Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(contents), &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
