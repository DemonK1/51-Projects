package models

import (
	"gorm.io/gorm"
	"log"
	"mysql_book_management_system/pkg/config"
)

var db *gorm.DB

// 初始化 mysql 链接等
func init() {
	config.Connect()
	db = config.GetDB()
	err := db.AutoMigrate(new(Book))
	if err != nil {
		log.Println(err)
		return
	}
}

// Book 用到的结构体
type Book struct {
	gorm.Model
	Name        string `gorm:"" json:"name"`
	Author      string `json:"author"`
	Publication string `json:"publication"`
}

func (b *Book) CreateBook() *Book {
	db.Create(&b)
	return b
}

func GetAllBook() []*Book {
	var Books []*Book
	db.Find(&Books)
	return Books
}

func GetBookById(Id int64) (*Book, *gorm.DB) {
	var getBook *Book
	db = db.Where("id=?", Id).Find(&getBook)
	return getBook, db
}

func DeleteBook(ID int64) *Book {
	var book Book
	db.Where("id=?", ID).Delete(&book)
	return &book
}
