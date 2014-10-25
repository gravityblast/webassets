package main

import (
	"fmt"
	"strings"
)

type stringslice []string

func (s *stringslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *stringslice) String() string {
	return fmt.Sprintf("[%s]", strings.Join(*s, ", "))
}
