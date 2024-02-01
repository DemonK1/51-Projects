package main

import (
	"crm_sqlite/database"
	"crm_sqlite/lead"
	"fmt"
	"github.com/gofiber/fiber/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Get("/api/v1/lead", lead.GetLeads)
	app.Post("/api/v1/lead", lead.PostLead)
	app.Delete("/api/v1/lead/:id", lead.DeleteLead)
}

func initDatabase() {
	var err error
	dsn := "sqlite3"
	database.DBConn, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("链接数据库出错")
	}
	fmt.Println("成功链接到数据库")
	l := new(lead.Lead)
	err = database.DBConn.AutoMigrate(l)
	if err != nil {
		return
	}
	fmt.Println("database Migrate")
}

func main() {
	app := fiber.New()
	setupRoutes(app)
	initDatabase()
	err := app.Listen("3000")
	if err != nil {
		panic(err)
	}

}
