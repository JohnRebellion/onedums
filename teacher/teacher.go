package teacher

import (
	"onedums/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// Teacher ...
type Teacher struct {
	gorm.Model `json:"-"`
	ID         uint      `json:"id" gorm:"primarykey"`
	UserID     uint      `json:"-" gorm:"unique"`
	User       user.User `json:"user"`
}

// GetTeachers ...
func GetTeachers(c *fiber.Ctx) error {
	teachers := []Teacher{}
	database.DBConn.Preload("User").Find(&teachers)
	return c.JSON(teachers)
}

// GetTeacher ...
func GetTeacher(c *fiber.Ctx) error {
	teacher := new(Teacher)
	database.DBConn.Preload("User").First(&teacher, c.Params("id"))
	return c.JSON(&teacher)
}

// NewTeacher ...
func NewTeacher(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	teacher := new(Teacher)
	fiberUtils.ParseBody(&teacher)
	database.DBConn.Create(&teacher)
	return fiberUtils.SendSuccessResponse("Created a new teacher successfully")
}

// UpdateTeacher ...
func UpdateTeacher(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	teacher := new(Teacher)
	fiberUtils.ParseBody(&teacher)
	database.DBConn.Updates(&teacher)
	return fiberUtils.SendSuccessResponse("Updated a teacher successfully")
}

// DeleteTeacher ...
func DeleteTeacher(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	database.DBConn.Delete(&Teacher{}, c.Params("id"))
	return fiberUtils.SendSuccessResponse("Deleted a teacher successfully")
}
