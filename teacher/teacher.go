package teacher

import (
	"onedums/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
	"github.com/JohnRebellion/go-utils/passwordHashing"
)

// Teacher ...
type Teacher struct {
	gorm.Model `json:"-"`
	ID         uint          `json:"id" gorm:"primarykey"`
	UserInfoID uint          `json:"-" gorm:"unique"`
	UserInfo   user.UserInfo `json:"userInfo" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// GetTeachers ...
func GetTeachers(c *fiber.Ctx) error {
	teachers := []Teacher{}
	teachersFiltered := []Teacher{}
	err := database.DBConn.Preload("UserInfo.User").Find(&teachers).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		for _, teacher := range teachers {
			if teacher.UserInfo.User.ID != 0 ||
				userClaim.User.Role == "Admin" {
				teachersFiltered = append(teachersFiltered, teacher)
			}
		}

		err = c.JSON(teachersFiltered)
	}

	return err
}

// GetTeacher ...
func GetTeacher(c *fiber.Ctx) error {
	teacher := new(Teacher)
	teacherFiltered := new(Teacher)
	err := database.DBConn.Preload("UserInfo.User").First(&teacher, c.Params("id")).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		if teacher.UserInfo.User.ID != 0 ||
			userClaim.User.Role == "Admin" {
			teacherFiltered = teacher
		}

		err = c.JSON(&teacherFiltered)
	}

	return err
}

// NewTeacher ...
func NewTeacher(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	teacher := new(Teacher)
	fiberUtils.ParseBody(&teacher)
	var err error
	teacher.UserInfo.User.Password, err = passwordHashing.HashPassword(teacher.UserInfo.User.Password)

	if err == nil {
		if database.DBConn.Create(&teacher).Error == nil {
			return fiberUtils.SendSuccessResponse("Created a new teacher successfully")
		}
	}

	return err
}

// UpdateTeacher ...
func UpdateTeacher(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	teacher := new(Teacher)
	fiberUtils.ParseBody(&teacher)
	var err error
	teacher.UserInfo.User.Password, err = passwordHashing.HashPassword(teacher.UserInfo.User.Password)

	if err == nil {
		if database.DBConn.Updates(&teacher).Error == nil {
			return fiberUtils.SendSuccessResponse("Updated a teacher successfully")
		}
	}

	return err
}

// DeleteTeacher ...
func DeleteTeacher(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&Teacher{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted a teacher successfully")
	}

	return err
}
