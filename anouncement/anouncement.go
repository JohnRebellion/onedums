package anouncement

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// Anouncement ...
type Anouncement struct {
	gorm.Model  `json:"-"`
	ID          uint      `json:"id" gorm:"primarykey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DateCreated time.Time `json:"dateCreated"`
}

// GetAnouncements ...
func GetAnouncements(c *fiber.Ctx) error {
	anouncements := []Anouncement{}
	err := database.DBConn.Find(&anouncements).Error

	if err == nil {
		err = c.JSON(anouncements)
	}

	return err
}

// GetAnouncement ...
func GetAnouncement(c *fiber.Ctx) error {
	anouncement := new(Anouncement)
	err := database.DBConn.First(&anouncement, c.Params("id")).Error

	if err == nil {
		err = c.JSON(&anouncement)
	}

	return err
}

// NewAnouncement ...
func NewAnouncement(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	anouncement := new(Anouncement)
	fiberUtils.ParseBody(&anouncement)
	anouncement.DateCreated = time.Now()
	err := database.DBConn.Create(&anouncement).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Created a new anouncement successfully")
	}

	return err
}

// UpdateAnouncement ...
func UpdateAnouncement(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	anouncement := new(Anouncement)
	fiberUtils.ParseBody(&anouncement)
	err := database.DBConn.Updates(&anouncement).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Updated a anouncement successfully")
	}

	return err
}

// DeleteAnouncement ...
func DeleteAnouncement(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&Anouncement{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted a anouncement successfully")
	}

	return err
}
