package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"reflect"
	"strconv"
	"text/template"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/seletskiy/tplutil"
)

type Config struct {
	Settings struct {
		//@TODO: move to separate structure and embed it here
		Carbohydrates float64 `toml:"carbohydrates"`
		Proteins      float64 `toml:"proteins"`
		Fats          float64 `toml:"fats"`
		Calories      float64 `toml:"calories"`
	}
}

const settingsShowTpl = `
	carbohydrates = {{.Settings.Carbohydrates}}{{"\n"}}
	proteins = {{.Settings.Proteins}}{{"\n"}}
	fats = {{.Settings.Fats}}{{"\n"}}
	calories = {{.Settings.Calories}}{{"\n"}}
`

func (config Config) String() string {
	myTpl := template.Must(
		template.New("settingsShowTpl").Parse(tplutil.Strip(
			settingsShowTpl,
		)))

	buf := bytes.NewBuffer([]byte{})
	myTpl.Execute(buf, config)

	return buf.String()
}

//@TODO: implement common function to handle file reading
func configRead(filename string) (Config, error) {
	config := Config{}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}

	if _, err := toml.Decode(string(contents), &config); err != nil {
		return Config{}, err
	}

	return config, nil
}

func configChange(config Config, entry string, value string) (Config, error) {
	//@TODO: implement changing by different config chunks, not
	//by hardcoded config.Settings
	field := reflect.
		Indirect(reflect.ValueOf(&config.Settings)).
		FieldByName(UCFirstLetter(entry))

	if !field.IsValid() {
		return config, errors.New("Not valid config entry")
	}

	switch field.Interface().(type) {
	case float64:
		floatValue, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return Config{}, err
		}
		field.SetFloat(floatValue)
	}

	return config, nil
}

func configWrite(filename string, config Config) error {
	buf := bytes.NewBuffer([]byte{})
	tomlEncoder := toml.NewEncoder(buf)
	tomlEncoder.Encode(config)

	err := ioutil.WriteFile(filename, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func UCFirstLetter(str string) string {
	processed := []rune(str)
	processed[0] = unicode.ToUpper(processed[0])
	return string(processed)
}
