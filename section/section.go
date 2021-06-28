package section

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// Section ...
type Section struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primarykey"`
	Name       string `json:"name"`
}

// GetSections ...
func GetSections(c *fiber.Ctx) error {
	sections := []Section{}
	database.DBConn.Find(&sections)
	return c.JSON(sections)
}

// GetSection ...
func GetSection(c *fiber.Ctx) error {
	section := new(Section)
	database.DBConn.First(&section, c.Params("id"))
	return c.JSON(&section)
}

// NewSection ...
func NewSection(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	section := new(Section)
	fiberUtils.ParseBody(&section)
	database.DBConn.Create(&section)
	return fiberUtils.SendSuccessResponse("Created a new section successfully")
}

// UpdateSection ...
func UpdateSection(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	section := new(Section)
	fiberUtils.ParseBody(&section)
	database.DBConn.Updates(&section)
	return fiberUtils.SendSuccessResponse("Updated a section successfully")
}

// DeleteSection ...
func DeleteSection(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	database.DBConn.Delete(&Section{}, c.Params("id"))
	return fiberUtils.SendSuccessResponse("Deleted a section successfully")
}
