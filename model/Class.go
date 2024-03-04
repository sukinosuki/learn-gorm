package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	Name     string
	Students []User `gorm:"foreignKey:class_id"` // Class拥有User, foreignKey:class_id对应User.ClassId字段
}
