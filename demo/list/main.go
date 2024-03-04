package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"learn-gorm-20240302/model"
	"time"
)

var Db *gorm.DB

var connectRetryCount = 0

func connectDb() {
	fmt.Println("连接mysql开始")
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details

	dsn := "root:233233@tcp(127.0.0.1:3306)/learn_gorm?charset=utf8mb4&parseTime=True&loc=Local"

	if connectRetryCount >= 3 {
		dsn = "root:233233@tcp(127.0.0.1:3306)/learn_gorm?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		//TranslateError: false,
	})

	if err != nil {
		fmt.Println("连接失败，准备重连， err: ", err.Error())
		time.Sleep(1 * time.Second)
		connectRetryCount++
		connectDb()
		return
	}
	fmt.Println("连接成功")
	connectRetryCount = 0
	Db = db.Debug()
}

func main() {
	fmt.Println("程序开始")
	//go connectDb()
	//time.Sleep(4 * time.Second)

	connectDb()
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

	fmt.Println("执行sql开始")
	//getUser(1)
	//associationGetUser(1)
	getUsers()
	//getClasses()
	//getBooks()
	//updateUser(1)
	//associationDelete(1)
	//createUser()
	//addClass()

	//getUsersByCertainFields()

	//testChainMethods()
	//testTransaction()
}

func testTransaction() error {

	tx := Db.Begin()

	defer func() {
		err := recover()
		if err != nil {
			tx.Rollback()
		}
	}()

	err := tx.Create(&model.UserIdCard{
		CardNumber: "222",
		Uid:        1,
	}).Error

	if err != nil {
		return err
	}

	err = tx.Create(&model.User{
		Name: "miiro",
	}).Error

	if err != nil {
		return err
	}

	return tx.Commit().Error
}

func testChainMethods() {
	start := time.Now()

	var users []model.User

	//tx := Db.Model(&model.User{})

	tx1 := Db.Where("name = ?", "hanami")
	tx2 := Db.Where("age < ?", 20)

	err := tx1.Find(&users).Error

	if err != nil {
		panic(err)
	}
	err = tx2.Find(&users).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("消耗ms: ", time.Now().UnixMilli()-start.UnixMilli())
}

func associationDelete(uid uint) {

	var user model.User
	user.ID = uid
	err := Db.Select("UserIdCard", "Books", "Languages").Delete(&user).Error
	if err != nil {
		panic(err)
	}
}

func createUser() {
	user := model.User{
		Name:    "hanami4",
		Age:     0,
		ClassId: 0,
	}
	// 有问题
	user.UserIdCard = &model.UserIdCard{
		CardNumber: "111",
	}

	err := Db.Create(&user).Error
	if err != nil {
		panic(err)
	}
}

func updateUser(uid uint) {

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

	var user model.User
	//user.ID = uid
	err := Db.Joins("UserIdCard").Joins("Class").Find(&user, uid).Error

	if err != nil {
		panic(err)
	}

	user.Name += "1"
	user.Age = 0
	err = Db.Model(&user).
		//Select("age", "name").
		Save(&user).Error

	if err != nil {
		panic(err)
	}
}

func associationGetUser(id uint) {
	var books []model.Book

	var user model.User
	user.ID = 1
	//err := Db.Find(&user, 1).Error
	//if err != nil {
	//	return
	//}
	err := Db.Model(&user).Where("author = ?", "aki").Association("Books").Find(&books)
	if err != nil {
		return
	}
	// Append: 往user(查出来的或者赋值id)->Books(一对多) 追加数据 (造成两条语句，一条新增book, 一条更新user.updated_at)
	// [rows:1] INSERT INTO `book` (`created_at`,`updated_at`,`deleted_at`,`name`,`author`,`uid`) VALUES ('2024-03-04 09:55:56.325','2024-03-
	//04 09:55:56.325',NULL,'book by append','author by append',1) ON DUPLICATE KEY UPDATE `uid`=VALUES(`uid`)

	// [rows:1] UPDATE `user` SET `updated_at`='2024-03-04 09:59:34.008' WHERE `user`.`deleted_at` IS NULL AND `id` = 2
	//err = Db.Model(&model.User{Model: gorm.Model{ID: 2}}).Association("Books").Append(&model.Book{Name: "book by append", Author: "author by append"})
	//err = Db.Model(&user).Association("Books").Append(&model.Book{Name: "book by append", Author: "author by append"})

	//  [rows:0] UPDATE `book` SET `uid`=NULL WHERE `book`.`id` <> 12 AND `book`.`uid` = 1 AND `book`.`deleted_at` IS NULL
	//  Error 1048 (23000): Column 'uid' cannot be null
	// Replace会将原来的books表对应id<>不等于新增 AND uid = 1的记录的 uid设置为null(不会删除记录), 如果book表uid字段设置为了不可为null则会抛出 Column 'uid' cannot be null 错误
	//err = Db.Model(&user).Association("Books").Replace(&model.Book{Name: "book by append", Author: "author by append"})

	// [rows:0] UPDATE `book` SET `uid`=NULL WHERE (`book`.`uid` = 1 AND `book`.`id` IN (1,2,6,8,9,10,11,12)) AND `book`.`deleted_at` IS NULL
	// Column 'uid' cannot be null
	//err = Db.Model(&user).Association("Books").Clear()

	//count := Db.Model(&user).Where("author = ?", "aki").Association("Books").Count()
	//fmt.Println("count ", count)

	err = Db.
		Unscoped(). // 逻辑删除
		Model(&user).Association("Books").
		Unscoped(). // 物理删除(原来是将book表外键uid设置为null, 加了unscoped会delete)
		Replace(&model.Book{Author: "author by unscoped replace", Name: "1111111"})

	if err != nil {
		return
	}
}
func getUser(uid uint) {
	//var users []model.User
	var user model.User
	err := Db.
		Joins("UserIdCard").
		Preload("Books").
		Where("user.id = ?", uid).
		First(&user).Error

	if err != nil {
		panic(err)
	}
	user.Name = "hanami"
	// 如果Books、UserIdCard有值, save方法会执行 [rows:0] insert into `user_id_card` ... ON DUPLICATE KEY UPDATE `uid`=VALUES(`uid`) ->结果是(books, user_id_card)数据不存在存新增
	// 推荐更新时使用select来指定要更新的字段来避免因为不熟悉gorm框架带来的不可预测的问题
	err = Db.Select("name").Save(user).Error

	if err != nil {
		panic(err)
	}

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

func getUsersByCertainFields() {

	var user []model.SimpleUser
	err := Db.Joins("Class").Preload("Books").Find(&user).Error
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

	// User-UserIdCard/Books/Class/Languages
	//err := Db.
	//	Joins("UserIdCard").
	//	Preload("Books", Db.Where("author = ?", "aki")).
	//	Joins("Class").
	//	Preload("Languages").
	//	//Where("user.name = ?", "miiro").
	//	Find(&users).Error

	// TODO: 查询指定字段
	err := Db.
		// 只对user.id 字段生效
		Select("user.id").
		//Joins("UserIdCard", func(db *gorm.DB) *gorm.DB {
		//	// Joins查询指定字段不生效
		//	// where不生效
		//	return db.Select("id").Where("card_number = ?", "1111")
		//}).
		// Joins查询指定字段，条件查询有效
		Joins("UserIdCard", Db.Select("id").Where("card_number = ?", "1111")).
		Preload("Books", func(db *gorm.DB) *gorm.DB {
			// preload查询一对多时，外键(这里是uid)是必需的
			return db.Select("uid")
		}).
		Joins("Class", func(db *gorm.DB) *gorm.DB {
			// Joins查询指定字段不生效
			//return db.Select("name")
			return db.Select("class.name")
		}).
		Preload("Languages", Db.Select("id")).
		//Preload("Languages", func(db *gorm.DB) *gorm.DB {
		//	// 查询Languages指定字段时，id是必需的
		//	return db.Select("id")
		//}).
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
