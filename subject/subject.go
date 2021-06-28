package subject

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// Subject ...
type Subject struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primarykey"`
	Name       string `json:"name"`
}

// GetSubjects ...
func GetSubjects(c *fiber.Ctx) error {
	subjects := []Subject{}
	database.DBConn.Find(&subjects)
	return c.JSON(subjects)
}

// GetSubject ...
func GetSubject(c *fiber.Ctx) error {
	subject := new(Subject)
	database.DBConn.First(&subject, c.Params("id"))
	return c.JSON(&subject)
}

// NewSubject ...
func NewSubject(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	subject := new(Subject)
	fiberUtils.ParseBody(&subject)
	database.DBConn.Create(&subject)
	return fiberUtils.SendSuccessResponse("Created a new subject successfully")
}

// UpdateSubject ...
func UpdateSubject(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	subject := new(Subject)
	fiberUtils.ParseBody(&subject)
	database.DBConn.Updates(&subject)
	return fiberUtils.SendSuccessResponse("Updated a subject successfully")
}

// DeleteSubject ...
func DeleteSubject(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	database.DBConn.Delete(&Subject{}, c.Params("id"))
	return fiberUtils.SendSuccessResponse("Deleted a subject successfully")
}
