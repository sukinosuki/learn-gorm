package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string
	Age        int
	Books      []Book      `gorm:"foreignKey:uid; references:id"`
	UserIdCard *UserIdCard `gorm:"foreignKey:uid"`
	ClassId    uint
	Class      Class      `gorm:"foreignKey:class_id"` // User属于Class, foreignKey:class_id对应User.ClassId字段
	Languages  []Language `gorm:"many2many:user_languages"`
}
