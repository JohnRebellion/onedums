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
	Role       string `json:"role" gorm:"default:user"`
	Status     string `json:"status"`
}

// GetUser Get a User by id
func GetUser(c *fiber.Ctx) error {
	user := new(User)
	database.DBConn.Find(&user, c.Params("id"))

	if user.ID == 0 {
		return c.JSON(fiber.Map{})
	}

	return c.JSON(user)
}

// GetUsers Get all users
func GetUsers(c *fiber.Ctx) error {
	users := []User{}
	database.DBConn.Preload("User").Find(&users)
	return c.JSON(users)
}

// NewUser Creating a User
func NewUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	user := new(User)
	err := fiberUtils.ParseBody(user)

	// return c.JSON(user)

	if err == nil {
		var usernameCount int64
		database.DBConn.Model(&user).Where("username = ?", user.Username).Count(&usernameCount)

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
			database.DBConn.Create(&user)
			return fiberUtils.SendJSONMessage("User Successfully Created", true, 201)
		}
	}

	return err
}

// UpdateUser Update User
func UpdateUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userO := new(User)
	err := fiberUtils.GetJWTClaimOfType("user", userO)

	if err == nil {
		user := new(User)
		err := fiberUtils.ParseBody(user)
		user.ID = user.ID

		if err == nil {
			if userO.Role == "Admin" || userO.ID == user.ID {
				if len(user.Username) == 0 || len(user.Password) == 0 || len(user.Name) == 0 {
					return fiberUtils.SendBadRequestResponse("Please Input Username, Password and Name")
				}

				if len(user.Username) < 3 || len(user.Password) < 8 || len(user.Name) < 3 {
					return fiberUtils.SendBadRequestResponse("Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively")
				}

				user.Password, err = passwordHashing.HashPassword(user.Password)
				fiberUtils.LogError(err)
				var existingUser User
				database.DBConn.First(&existingUser, user.ID)

				if existingUser.ID == 0 {
					return fiberUtils.SendJSONMessage("No User exists", false, 404)
				}

				database.DBConn.Updates(&user)
				return fiberUtils.SendSuccessResponse("User Successfully Updated")
			}

			err = fiberUtils.SendJSONMessage("No permission to update", false, 401)
		}
	}

	return err
}

// DeleteUser Delete User by id
func DeleteUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userO := new(User)
	err := fiberUtils.GetJWTClaimOfType("user", userO)

	if err == nil {
		user := new(User)
		database.DBConn.First(&user, c.Params("id"))

		if user.ID == 0 {
			return fiberUtils.SendJSONMessage("No User exists", false, 404)
		}

		if userO.Role == "Admin" || userO.ID == user.ID {
			database.DBConn.Delete(&user)
			return fiberUtils.SendSuccessResponse("User Successfully Deleted")
		}

		err = fiberUtils.SendJSONMessage("No permission to delete", false, 401)
	}

	return err
}

// GetUserFromJWTClaim ...
func GetUserFromJWTClaim(c *fiber.Ctx) User {
	fiberUtils.Ctx.New(c)
	userO := fiberUtils.GetJWTClaim("user")
	user := new(User)
	database.DBConn.First(&user.ID, userO["id"])
	return *user
}
