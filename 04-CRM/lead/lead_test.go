package lead

import (
	"fmt"
	"github.com/gofiber/fiber/v3"
	"testing"
)

func TestGetLead(t *testing.T) {
	//app := fiber.New()
	var c fiber.Ctx
	err := GetLead(c)
	if err != nil {
		fmt.Println(err)
	}
}
