package main

import (
	"io"
	"io/ioutil"

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
