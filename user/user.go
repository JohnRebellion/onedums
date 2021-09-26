// Package user provides user database table and CRUD API
package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
	"github.com/JohnRebellion/go-utils/passwordHashing"
)

// User ...
type User struct {
	gorm.Model `json:"-"`
	ID         uint   `json:"id" gorm:"primarykey"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Role       string `json:"role" gorm:"default:User"`
	Status     string `json:"status"`
}

// GetUser Get a User by id
func GetUser(c *fiber.Ctx) error {
	user := new(User)
	err := database.DBConn.Find(&user, c.Params("id")).Error

	if err == nil {
		err = c.JSON(user)
	}

	return err
}

// GetUsers Get all users
func GetUsers(c *fiber.Ctx) error {
	users := []User{}
	err := database.DBConn.Preload("User").Find(&users).Error

	if err == nil {
		err = c.JSON(users)
	}

	return err
}

// NewUser Creating a User
func NewUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	user := new(User)
	err := fiberUtils.ParseBody(user)

	if err == nil {
		var usernameCount int64

		if database.DBConn.Model(&user).Where("username = ?", user.Username).Count(&usernameCount).Error == nil {
			if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Name) == 0 {
				return fiberUtils.SendBadRequestResponse("Please Input Username, Password and Name")
			}

			if len(user.Username) < 3 || len(user.Password) < 8 || len(user.Name) < 3 {
				return fiberUtils.SendBadRequestResponse("Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively")
			}

			if usernameCount > 0 {
				return fiberUtils.SendBadRequestResponse("Username Already Exists")
			}

			user.Password, err = passwordHashing.HashPassword(user.Password)

			if err == nil {
				if database.DBConn.Create(&user).Error == nil {
					return fiberUtils.SendJSONMessage("User Successfully Created", true, 201)
				}
			}
		}
	}

	return err
}

// UpdateUser Update User
func UpdateUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	user := new(User)
	err := fiberUtils.ParseBody(user)
	userClaim := GetUserInfoFromJWTClaim(c)

	if err == nil {
		if userClaim.User.Role == "Admin" || userClaim.User.ID == user.ID {
			if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Name) == 0 {
				return fiberUtils.SendBadRequestResponse("Please Input Username, Password and Name")
			}

			if len(user.Username) < 3 || len(user.Password) < 8 || len(user.Name) < 3 {
				return fiberUtils.SendBadRequestResponse("Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively")
			}

			user.Password, err = passwordHashing.HashPassword(user.Password)
			fiberUtils.LogError(err)
			var existingUser User

			if database.DBConn.First(&existingUser, user.ID).Error == nil {
				if existingUser.ID == 0 {
					return fiberUtils.SendJSONMessage("No User exists", false, 404)
				}

				if database.DBConn.Updates(&user).Error == nil {
					return fiberUtils.SendSuccessResponse("User Successfully Updated")
				}
			}
		}

		return fiberUtils.SendJSONMessage("No permission to update", false, 401)
	}

	return err
}

// DeleteUser Delete User by id
func DeleteUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userClaim := GetUserInfoFromJWTClaim(c)
	user := new(User)
	err := database.DBConn.First(&user, c.Params("id")).Error

	if err == nil {
		if user.ID == 0 {
			return fiberUtils.SendJSONMessage("No User exists", false, 404)
		}

		if userClaim.User.Role == "Admin" || userClaim.User.ID == user.ID {
			if database.DBConn.Delete(&user).Error == nil {
				return fiberUtils.SendSuccessResponse("User Successfully Deleted")
			}
		}

		return fiberUtils.SendJSONMessage("No permission to delete", false, 401)
	}

	return err
}
