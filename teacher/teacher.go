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
	teacher.UserInfo.User.Role = "Teacher"
	teacher.UserInfo.User.Password, err = passwordHashing.HashPassword(teacher.UserInfo.User.Password)
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		if userClaim.User.ID == teacher.UserInfo.User.ID || userClaim.User.Role == "Admin" {
			if database.DBConn.Create(&teacher).Error == nil {
				return fiberUtils.SendSuccessResponse("Created a new teacher successfully")
			}
		} else {
			return fiberUtils.SendJSONMessage("No permission to delete", false, 401)
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
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		if userClaim.User.ID == teacher.UserInfo.User.ID || userClaim.User.Role == "Admin" {
			if database.DBConn.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&teacher).Error == nil {
				return fiberUtils.SendSuccessResponse("Updated a teacher successfully")
			}
		} else {
			return fiberUtils.SendJSONMessage("No permission to delete", false, 401)
		}
	}

	return err
}

// DeleteTeacher ...
func DeleteTeacher(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userClaim := user.GetUserInfoFromJWTClaim(c)
	teacher := new(Teacher)
	err := database.DBConn.First(&teacher, c.Params("id")).Error

	if err == nil {
		if userClaim.User.ID == teacher.UserInfo.User.ID || userClaim.User.Role == "Admin" {
			if database.DBConn.Delete(&teacher).Error == nil {
				return fiberUtils.SendSuccessResponse("Updated a teacher successfully")
			}
		} else {
			return fiberUtils.SendJSONMessage("No permission to delete", false, 401)
		}
	}

	return err
}
