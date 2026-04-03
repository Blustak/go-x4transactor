package config

import (
	"bufio"
	"strings"
)

type Config struct{
    bufio.Writer
    bufio.Scanner
}

func (c *Config) GetArgs() []string {
    if c.Scan() {
        return strings.Fields(c.Text())
    } else {
        return nil
    }
}
