package main

import (
	"flag"
	"fmt"
	"strings"
)

var _ flag.Value = &prefixMapFlag{}

type prefixMapFlag map[string]string

func (p *prefixMapFlag) String() string {
	return fmt.Sprint(*p)
}

func (p *prefixMapFlag) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid format, expected key:value")
	}
	(*p)[parts[0]] = parts[1]
	return nil
}
