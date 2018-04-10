package main

import (
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model

	UserID uint   `gorm:"not null"`
	Token  string `gorm:"size:36;not null;unique"`
}

type PasteRecord struct {
	gorm.Model

	UserID  uint   `gorm:"not null;index"`
	Content string `gorm:"not null" sql:"type:text"`
}
