package student

import (
	"onedums/section"
	"onedums/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// Student ...
type Student struct {
	gorm.Model `json:"-"`
	ID         uint            `json:"id" gorm:"primarykey"`
	UserID     uint            `json:"-" gorm:"unique"`
	User       user.User       `json:"user"`
	SectionID  uint            `json:"-"`
	Section    section.Section `json:"section"`
}

// GetStudents ...
func GetStudents(c *fiber.Ctx) error {
	students := []Student{}
	database.DBConn.Preload("Section").Preload("User").Find(&students)
	return c.JSON(students)
}

// GetStudent ...
func GetStudent(c *fiber.Ctx) error {
	student := new(Student)
	database.DBConn.Preload("Section").Preload("User").First(&student, c.Params("id"))
	return c.JSON(&student)
}

// NewStudent ...
func NewStudent(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	student := new(Student)
	fiberUtils.ParseBody(&student)
	database.DBConn.Create(&student)
	return fiberUtils.SendSuccessResponse("Created a new student successfully")
}

// UpdateStudent ...
func UpdateStudent(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	student := new(Student)
	fiberUtils.ParseBody(&student)
	database.DBConn.Updates(&student)
	return fiberUtils.SendSuccessResponse("Updated a student successfully")
}

// DeleteStudent ...
func DeleteStudent(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	database.DBConn.Delete(&Student{}, c.Params("id"))
	return fiberUtils.SendSuccessResponse("Deleted a student successfully")
}
