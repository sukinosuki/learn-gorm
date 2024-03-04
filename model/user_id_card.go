package model

import "gorm.io/gorm"

type UserIdCard struct {
	gorm.Model
	Uid        uint
	CardNumber string
}
