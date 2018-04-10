package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config is a global configuration struct
type Config struct {
	RPCAddr string `json:"rpc_addr"`
	DBPath  string `json:"db_path"`
	Token   string `json:"token"`
}

// NewConfig return a instance of Config
func NewConfig() *Config {
	c := &Config{}
	// read from config file
	jsonBytes, err := ioutil.ReadFile(os.Getenv("HOME") + "/pasty.json")
	if err != nil {
		log.Printf("failed to load config file: %s", err)
	} else {
		if err := json.Unmarshal(jsonBytes, c); err != nil {
			log.Printf("failed to load config file: %s", err)
		}
	}

	// if dbPath is not set, set default
	if c.DBPath == "" {
		c.DBPath = "./pasty.db"
	}

	// if dbPath flag is on, use it
	if *dbPath != "" {
		c.DBPath = *dbPath
	}

	// if rpc addr is not set, set default
	if c.RPCAddr == "" {
		c.RPCAddr = "127.0.0.1:9527"
	}

	// if rpcAddr flag is on, use it
	if *rpcAddr != "" {
		c.RPCAddr = *rpcAddr
	}

	return c
}
