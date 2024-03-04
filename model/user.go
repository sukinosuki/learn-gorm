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
	Class      Class `gorm:"foreignKey:class_id"` // User属于Class, foreignKey:class_id对应User.ClassId字段
	// User拥有多个Language, 使用many2many:[外键表], 这里的外键表为user_languages
	// 外键表user_languages字段为user_id(User.id), language_id(Language.id)(也可以自定义)
	Languages []Language `gorm:"many2many:user_languages"`
}

type SimpleUser struct {
	ID      uint
	Name    string
	ClassId uint
	Class   SimpleClass
	Books   []SimpleBook `gorm:"foreignKey:uid"`
}
type SimpleBook struct {
	ID   uint
	Name string
	UID  uint
}

func (b *SimpleBook) TableName() string {
	return "book"
}

func (u *SimpleUser) TableName() string {
	return "user"
}

type SimpleClass struct {
	ID   uint
	Name string
}

func (c *SimpleClass) TableName() string {
	return "class"
}
