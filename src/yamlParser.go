// Package gotapper    provides ...
package gotapper

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func readConfigFile(fileName string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	} else {
		return bytes, nil
	}
}

func parse(fileName string) (config, error) {
	fileBytes, err := readConfigFile(fileName)
	if err != nil {
		return config{}, err
	}
	configFromFile := config{}
	if err = yaml.Unmarshal(fileBytes, &configFromFile); err != nil {
		return configFromFile, nil
	}
	return configFromFile, err
}
