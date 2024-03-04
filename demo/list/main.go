package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"learn-gorm-20240302/model"
)

var Db *gorm.DB

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:233233@tcp(127.0.0.1:3306)/learn_gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	Db = db.Debug()

	if err != nil {
		panic(err)
	}

	// 1-----
	//var books []map[string]any
	//err = Db.Model(&model.Book{}).Select("user.*, book.*").Joins("LEFT JOIN user ON user.id = book.uid").Find(&books).Error

	// 1-----

	// 2-----
	//var res []map[string]any
	//bookQuery := db.Model(&model.Book{}).Where("author = ?", "aki")
	//err = Db.Model(&model.User{}).
	//	Select("user.*, book.id bookId, book.author, book.uid bookUid").
	//	Joins("LEFT JOIN (?) book ON Book.uid = user.id", bookQuery).
	//	Find(&res).Error
	//
	//if err != nil {
	//	panic(err)
	//}
	// 2-----

	//addUser(db)
	//addBook(2)
	//addIdCard(2)

	//getUser(1)
	getUsers()
	//getClasses()
	//getBooks()
	//updateUser(1)

	//addClass()
}

func updateUser(uid int) {

	// 1---
	//var user model.User
	//err := Db.Model(&model.User{}).Where("id = ? ", uid).First(&user).Error
	//if err != nil {
	//	panic(err)
	//}
	//user.Name = user.Name + "1"
	//user.ClassId = 0 // save会更新0值
	//err = Db.Save(&user).Error
	//if err != nil {
	//	panic(err)
	//}
	// 1---

	//err := Db.Table("user").Where("id = ?", uid).Update("name", "0").Update("age", 0).Error

	err := Db.Model(&model.User{}).
		Select("age", "name").
		Where("id = ?", uid).
		Updates(&model.User{Age: 0, Name: ""}).Error

	if err != nil {
		panic(err)
	}
}

func getUser(uid uint) {
	//var users []model.User
	var user model.User
	err := Db.Model(&model.User{}).Preload("UserIdCard").Where("id = ?", uid).First(&user).Error

	if err != nil {
		panic(err)
	}

	fmt.Println("user ", user)
}

func getBooks() {
	var books []model.Book

	err := Db.
		Joins("User").
		Joins("User.UserIdCard").
		Joins("User.Class").
		Find(&books).
		Preload("User.Books").
		Preload("User.Languages").
		Find(&books).
		Error

	if err != nil {
		panic(err)
	}
}

func getClasses() {

	var list []model.Class

	err := Db.
		Preload("Students").
		Preload("Students.Class").
		Preload("Students.Books", func(db *gorm.DB) *gorm.DB {
			return db.Limit(1)
		}).
		Preload("Students.Languages").
		Find(&list).Error

	if err != nil {
		panic(err)
	}
}

func getUsers() {
	var users []model.User
	//var user model.User
	//err := Db.Model(&model.User{}).Preload("UserIdCard").Preload("Books2").Find(&users).Error

	//err := Db.
	//	Joins("left join class on `user`.class_id = class.id").
	//	Where("user.id =?", 1).
	//	Find(&users).Error

	err := Db.
		//Joins("UserIdCard", Db.Where("card_number = ?", "1112")).
		Joins("UserIdCard").
		Preload("Books", Db.Where("author = ?", "aki")).
		Joins("Class").
		Preload("Languages").
		//Where("user.name = ?", "miiro").
		Find(&users).Error

	// preload 条件查询
	//err := Db.
	//	Preload("Class.Students", "id = ?", 1).
	//	Preload(clause.Associations).
	//	Limit(1).
	//	Find(&users).Error

	// preload 其它自定义操作
	//err := Db.
	//	Preload("Class.Students", func(db *gorm.DB) *gorm.DB {
	//		// 可以条件查询
	//		//return db.Where("id = ?", 1)
	//
	//		return db.Order("id DESC")
	//	}).
	//	Preload(clause.Associations).
	//	Limit(10).
	//	Find(&users).Error

	if err != nil {
		panic(err)
	}

}

func addClass() {
	class := &model.Class{
		Name: "class2",
	}
	err := Db.Save(&class).Error
	if err != nil {
		panic(err)
	}
}

func addUser(db *gorm.DB) {
	user := model.User{
		Name: "miiro",
	}

	err := db.Save(&user).Error
	if err != nil {
		panic(err)
	}
}

func addBook(uid uint) {
	book := &model.Book{
		Name:   "book444",
		Author: "author444",
		Uid:    uid,
	}
	err := Db.Save(&book).Error
	if err != nil {
		panic(err)
	}
}

func addIdCard(uid uint) {
	idCard := &model.UserIdCard{
		Uid:        uid,
		CardNumber: "2222",
	}
	err := Db.Save(&idCard).Error
	if err != nil {
		panic(err)
	}
}
