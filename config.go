package main

import (
	"os"
)

// Config is a global configuration struct
type Config struct {
	Addr string
}

// NewConfig return a instance of Config
func NewConfig() *Config {
	c := &Config{}
	c.Addr = os.Getenv("rpcAddr")
	if c.Addr == "" {
		c.Addr = "127.0.0.1:9527"
	}

	return c
}
