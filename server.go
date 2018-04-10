package main

import (
	"context"
	"errors"
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
	var items []*PasteItem

	rows, err := db.Table(
		"paste_record",
	).Select(
		"paste_record.id, paste_record.updated_at, paste_record.content",
	).Joins(
		"LEFT JOIN token ON paste_record.user_id = token.user_id",
	).Where(
		"token.token = ?", req.Token,
	).Order(
		"paste_record.id DESC",
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
	var token Token
	db.Where("token = ?", req.Token).Find(&token)
	if token.UserID == 0 {
		return nil, errors.New("token not valid, please login first")
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
