package main

import (
	"flag"
	"os"

	"github.com/jinzhu/gorm"
)

var (
	isServer = flag.Bool("isServer", false, "isServer is set if it's a server process")
	dbPath   = flag.String("dbPath", "", "default db path")
	rpcAddr  = flag.String("rpcAddr", "", "rpc address")
	limit    = flag.Int64("limit", 0, "how many records do we need, at most 100")
	hint     = flag.Bool("hint", true, "enable hint or disable")

	config *Config
	db     *gorm.DB
)

func init() {
	flag.Parse()

	config = NewConfig()

	// make sure db file exist
	if *isServer {
		if _, err := os.Stat(config.DBPath); os.IsNotExist(err) {
			os.OpenFile(config.DBPath, os.O_RDONLY|os.O_CREATE, 0700)
		}
	}
}

func main() {
	if *isServer {
		startPastyServer()
	}
	startPastyClient()
}
