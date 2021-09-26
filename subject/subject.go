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
	err := database.DBConn.Find(&subjects).Error

	if err == nil {
		err = c.JSON(subjects)
	}

	return err
}

// GetSubject ...
func GetSubject(c *fiber.Ctx) error {
	subject := new(Subject)
	err := database.DBConn.First(&subject, c.Params("id")).Error

	if err == nil {
		err = c.JSON(&subject)
	}

	return err
}

// NewSubject ...
func NewSubject(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	subject := new(Subject)
	fiberUtils.ParseBody(&subject)
	err := database.DBConn.Create(&subject).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Created a new subject successfully")
	}

	return err
}

// UpdateSubject ...
func UpdateSubject(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	subject := new(Subject)
	fiberUtils.ParseBody(&subject)
	err := database.DBConn.Updates(&subject).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Updated a subject successfully")
	}

	return err
}

// DeleteSubject ...
func DeleteSubject(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&Subject{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted a subject successfully")
	}

	return err
}
