package model

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	gorm.Model
	Name      string
	Author    string
	Uid       uint
	User      User `gorm:"foreignKey:uid"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
