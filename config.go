package main

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Javascripts map[string][]string
}

func parseConfig(r io.Reader) (config, error) {
	c := config{}
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(content, &c)

	return c, err
}

func parseConfigFile(path string) (config, error) {
	f, err := os.Open(path)
	if err != nil {
		return config{}, err
	}

	return parseConfig(f)
}
