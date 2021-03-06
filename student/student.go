package student

import (
	"onedums/section"
	"onedums/user"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
	"github.com/JohnRebellion/go-utils/passwordHashing"
)

// Student ...
type Student struct {
	gorm.Model `json:"-"`
	ID         uint            `json:"id" gorm:"primarykey"`
	UserInfoID uint            `json:"-" gorm:"unique"`
	UserInfo   user.UserInfo   `json:"userInfo"`
	SectionID  uint            `json:"-"`
	Section    section.Section `json:"section"`
	Guardian   Guardian        `json:"guardian" gorm:"embedded"`
}

// Guardian
type Guardian struct {
	Name          string `json:"name"`
	ContactNumber string `json:"contactNumber"`
}

// GetStudents ...
func GetStudents(c *fiber.Ctx) error {
	students := []Student{}
	studentsFiltered := []Student{}
	err := database.DBConn.Preload("Section").Preload("UserInfo.User").Find(&students).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		for _, student := range students {
			if student.Section.ID != 0 &&
				student.UserInfo.User.ID != 0 ||
				userClaim.User.Role == "Admin" {
				studentsFiltered = append(studentsFiltered, student)
			}
		}

		err = c.JSON(studentsFiltered)
	}

	return err
}

// GetStudent ...
func GetStudent(c *fiber.Ctx) error {
	student := new(Student)
	studentFiltered := new(Student)
	err := database.DBConn.Preload("Section").Preload("UserInfo.User").First(&student, c.Params("id")).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		if student.Section.ID != 0 &&
			student.UserInfo.User.ID != 0 ||
			userClaim.User.Role == "Admin" {
			studentFiltered = student
		}

		err = c.JSON(&studentFiltered)
	}

	return err
}

// NewStudent ...
func NewStudent(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	student := new(Student)
	err := fiberUtils.ParseBody(&student)

	if err == nil {
		student.UserInfo.User.Role = "Student"
		student.UserInfo.User.Password, err = passwordHashing.HashPassword(student.UserInfo.User.Password)

		if err == nil {
			if database.DBConn.Session(&gorm.Session{FullSaveAssociations: true}).Create(&student).Error == nil {
				return fiberUtils.SendSuccessResponse("Created a new student successfully")
			}
		}
	}

	return err
}

// UpdateStudent ...
func UpdateStudent(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	student := new(Student)
	err := fiberUtils.ParseBody(&student)

	if err == nil {
		student.UserInfo.User.Role = "Student"
		student.UserInfo.User.Password, err = passwordHashing.HashPassword(student.UserInfo.User.Password)
		userClaim := user.GetUserInfoFromJWTClaim(c)

		if err == nil {
			if userClaim.User.ID == student.UserInfo.User.ID || userClaim.User.Role == "Admin" {
				if database.DBConn.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&student).Error == nil {
					return fiberUtils.SendSuccessResponse("Updated a student successfully")
				}
			} else {
				return fiberUtils.SendJSONMessage("No permission to delete", false, 401)
			}
		}
	}

	return err
}

// DeleteStudent ...
func DeleteStudent(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userClaim := user.GetUserInfoFromJWTClaim(c)
	student := new(Student)
	err := database.DBConn.First(&student, c.Params("id")).Error

	if err == nil {
		if userClaim.User.ID == student.UserInfo.User.ID || userClaim.User.Role == "Admin" {
			if database.DBConn.Session(&gorm.Session{FullSaveAssociations: true}).Delete(&student).Error == nil {
				return fiberUtils.SendSuccessResponse("Updated a student successfully")
			}
		} else {
			return fiberUtils.SendJSONMessage("No permission to delete", false, 401)
		}
	}

	return err
}
