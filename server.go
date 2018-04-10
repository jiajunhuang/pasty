package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type server struct{}

func (s *server) GetPaste(ctx context.Context, req *GetPasteRequest) (*PasteResponse, error) {
	log.Println("GetPaste been called")
	var items []*PasteItem

	rows, err := db.Table(
		"paste_records",
	).Select(
		"paste_records.id, paste_records.updated_at, paste_records.content",
	).Joins(
		"LEFT JOIN tokens ON paste_records.user_id = tokens.user_id",
	).Where(
		"tokens.token = ?", req.Token,
	).Order(
		"paste_records.id DESC",
	).Limit(
		req.Limit,
	).Rows()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var recordID uint
		var updatedAt time.Time
		var content string

		if err := rows.Scan(&recordID, &updatedAt, &content); err == nil {
			items = append(items, &PasteItem{Id: int64(recordID), Timestamp: updatedAt.Unix(), Content: content})
		}
	}
	return &PasteResponse{Items: items}, nil
}

func (s *server) Paste(ctx context.Context, req *PasteRequest) (*PasteResponse, error) {
	log.Println("Paste been called")
	var token Token
	db.Where("token = ?", req.Token).Find(&token)
	if token.UserID == 0 {
		return nil, fmt.Errorf("token %s not valid, please login first", req.Token)
	}

	record := PasteRecord{
		UserID:  token.UserID,
		Content: req.Content,
	}
	db.Create(&record)
	return &PasteResponse{Items: []*PasteItem{}}, nil
}

func startPastyServer() {
	var err error
	log.Printf("gonna open db: %s", config.DBPath)
	db, err = gorm.Open("sqlite3", config.DBPath)
	if err != nil {
		log.Panicf("failed to open db %s: %s", config.DBPath, err)
	}
	db.AutoMigrate(&Token{}, &PasteRecord{})

	lis, err := net.Listen("tcp", config.RPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("listen at: %v", config.RPCAddr)
	s := grpc.NewServer()
	RegisterPastyServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
