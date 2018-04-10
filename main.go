package main

import (
	"flag"
	"os"

	"github.com/jinzhu/gorm"
)

var (
	isServer  = flag.Bool("isServer", false, "isServer is set if it's a server process")
	configDir = flag.String("configDir", "~/.config/pasty/", "default config dir")
	dbPath    = flag.String("dbPath", "", "default db path")
	rpcAddr   = flag.String("rpcAddr", "", "rpc address")

	config *Config
	db     *gorm.DB
)

func init() {
	config = NewConfig()

	// make sure directory exist
	if _, err := os.Stat(*configDir); os.IsNotExist(err) {
		os.Mkdir(*configDir+"pasty.json", 0700)
	}

	// make sure db file exist
	if _, err := os.Stat(config.DBPath); os.IsNotExist(err) {
		os.OpenFile(config.DBPath, os.O_RDONLY|os.O_CREATE, 0700)
	}
}

func main() {
	flag.Parse()

	if *isServer {
		startPastyServer()
	}
	startPastyClient()
}
