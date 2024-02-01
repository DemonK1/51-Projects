package lead

import (
	"crm_sqlite/database"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"log"
)

type Lead struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
	Phone   int    `json:"phone"`
	// gorm 自带的一些数据: ID CreatedAt UpdatedAt DeletedAt 我把他们放在最底部
	gorm.Model
}

func GetLead(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	db := database.DBConn
	lead := new(Lead)
	db.Find(lead, id)
	err := ctx.JSON(lead)
	if err != nil {
		log.Println("响应 json 解码错误: ", err)
	}
	return nil
}

func GetLeads(ctx fiber.Ctx) error {
	return nil
}

func PostLead(ctx fiber.Ctx) error {
	return nil
}

func DeleteLead(ctx fiber.Ctx) error {
	return nil
}
