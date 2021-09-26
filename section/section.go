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
	err := database.DBConn.Find(&sections).Error

	if err == nil {
		err = c.JSON(sections)
	}

	return err
}

// GetSection ...
func GetSection(c *fiber.Ctx) error {
	section := new(Section)
	err := database.DBConn.First(&section, c.Params("id")).Error

	if err == nil {
		err = c.JSON(&section)
	}

	return err
}

// NewSection ...
func NewSection(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	section := new(Section)
	fiberUtils.ParseBody(&section)
	err := database.DBConn.Create(&section).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Created a new section successfully")
	}

	return err
}

// UpdateSection ...
func UpdateSection(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	section := new(Section)
	fiberUtils.ParseBody(&section)
	err := database.DBConn.Updates(&section).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Updated a section successfully")
	}

	return err
}

// DeleteSection ...
func DeleteSection(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&Section{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted a section successfully")
	}

	return err
}
