package main

import (
	"github.com/gorilla/mux"
	"log"
	"mysql_book_management_system/pkg/routers"
	"net/http"
)

// var db *gorm.DB

// // 初始化 mysql 链接等
// func init() {
// 	config.Connect()
// 	db = config.GetDB()
// 	err := db.AutoMigrate(new(models.Book))
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }

func main() {
	r := mux.NewRouter()
	routers.RegisterBookStoreRoutes(r)
	log.Fatal(http.ListenAndServe(":8080", r))
}
