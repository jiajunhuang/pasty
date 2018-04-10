package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func startPastyClient() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(config.RPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := NewPastyClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	limitV := *limit
	if limitV == 0 {
		// read from stdin and create paste
		content, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Printf("failed to read: %s", err)
			return
		}

		if _, err := c.Paste(ctx, &PasteRequest{Token: config.Token, Content: string(content)}); err != nil {
			log.Printf("failed to paste: %s", err)
		}
		return
	}

	// read pastes from server
	if limitV < 0 {
		limitV = 1
	}
	if limitV > 100 {
		limitV = 100
	}

	// Contact the server and print out its response.
	r, err := c.GetPaste(ctx, &GetPasteRequest{Token: config.Token, Limit: limitV})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if len(r.Items) == 0 {
		return
	}

	for _, item := range r.Items {
		if *hint {
			fmt.Printf("file id %d created at %v      ==================\n", item.Id, time.Unix(item.Timestamp, 0))
		}
		fmt.Println(item.Content)
		if *hint {
			fmt.Printf("file id %d created at %v over ==================\n", item.Id, time.Unix(item.Timestamp, 0))
		}
	}
}
