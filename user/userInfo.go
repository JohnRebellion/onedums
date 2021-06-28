// Package user provides user database table and CRUD API
package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
	"github.com/JohnRebellion/go-utils/passwordHashing"
)

type UserInfo struct {
	gorm.Model    `json:"-"`
	ID            uint      `json:"id" gorm:"primarykey"`
	UserID        uint      `json:"-" gorm:"unique"`
	User          User      `json:"user"`
	Birthday      time.Time `json:"birthday"`
	Gender        string    `json:"gender"`
	Address       string    `json:"address"`
	ContactNumber string    `json:"contactNumber"`
}

// GetUserInfo Get a User by id
func GetUserInfo(c *fiber.Ctx) error {
	userInfo := new(UserInfo)
	database.DBConn.Preload("User").Find(&userInfo, c.Params("id"))

	if userInfo.ID == 0 {
		return c.JSON(fiber.Map{})
	}

	return c.JSON(userInfo)
}

// GetUserInfos Get all users
func GetUserInfos(c *fiber.Ctx) error {
	userInfos := []UserInfo{}
	database.DBConn.Preload("User").Find(&userInfos)
	return c.JSON(userInfos)
}

// NewUserInfo Creating a User
func NewUserInfo(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userInfo := new(UserInfo)
	err := fiberUtils.ParseBody(userInfo)

	// return c.JSON(userInfo)

	if err == nil {
		var usernameCount int64
		database.DBConn.Model(&userInfo.User).Where("username = ?", userInfo.User.Username).Count(&usernameCount)

		if len(userInfo.User.Username) == 0 || len(userInfo.User.Password) == 0 || len(userInfo.User.Name) == 0 {
			return fiberUtils.SendBadRequestResponse("Please Input Username, Password and Name")
		}

		if len(userInfo.User.Username) < 3 || len(userInfo.User.Password) < 8 || len(userInfo.User.Name) < 3 {
			return fiberUtils.SendBadRequestResponse("Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively")
		}

		if usernameCount > 0 {
			return fiberUtils.SendBadRequestResponse("Username Already Exists")
		}

		userInfo.User.Password, err = passwordHashing.HashPassword(userInfo.User.Password)

		if err == nil {
			database.DBConn.Create(&userInfo)
			return fiberUtils.SendJSONMessage("User Successfully Created", true, 201)
		}
	}

	return err
}

// Authenticate the User and return User object and token string with cookies
func Authenticate(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	user := new(User)
	err := fiberUtils.ParseBody(user)
	password := user.Password

	if err == nil {
		if len(user.Username) == 0 || len(password) == 0 {
			return fiberUtils.SendBadRequestResponse("Please Input Username and Password")
		}

		userInfo := UserInfo{}
		database.DBConn.Joins("User").First(&userInfo, "username = ?", user.Username)

		if passwordHashing.CheckPasswordHash(password, userInfo.User.Password) {
			_, err := fiberUtils.GenerateJWTSignedString(fiber.Map{"userInfo": userInfo, "isAdmin": userInfo.User.Role == "Admin"})

			if err == nil {
				return fiberUtils.SendSuccessResponse("Access Granted!")
			}
		} else {
			return fiberUtils.SendJSONMessage("Incorrect Username or Password", false, 401)
		}
	}

	return err
}

// UpdateUserInfo Update User
func UpdateUserInfo(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userO := new(UserInfo)
	err := fiberUtils.GetJWTClaimOfType("userInfo", userO)

	if err == nil {
		userInfo := new(UserInfo)
		err := fiberUtils.ParseBody(userInfo)
		userInfo.User.ID = userInfo.ID

		if err == nil {
			if userO.User.Role == "Admin" || userO.ID == userInfo.ID {
				if len(userInfo.User.Username) == 0 || len(userInfo.User.Password) == 0 || len(userInfo.User.Name) == 0 {
					return fiberUtils.SendBadRequestResponse("Please Input Username, Password and Name")
				}

				if len(userInfo.User.Username) < 3 || len(userInfo.User.Password) < 8 || len(userInfo.User.Name) < 3 {
					return fiberUtils.SendBadRequestResponse("Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively")
				}

				userInfo.User.Password, err = passwordHashing.HashPassword(userInfo.User.Password)
				fiberUtils.LogError(err)
				var existingUserInfo UserInfo
				database.DBConn.First(&existingUserInfo, userInfo.ID)

				if existingUserInfo.ID == 0 {
					return fiberUtils.SendJSONMessage("No User exists", false, 404)
				}

				database.DBConn.Updates(&userInfo)
				return fiberUtils.SendSuccessResponse("User Successfully Updated")
			}

			err = fiberUtils.SendJSONMessage("No permission to update", false, 401)
		}
	}

	return err
}

// DeleteUserInfo Delete User by id
func DeleteUserInfo(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userO := new(UserInfo)
	err := fiberUtils.GetJWTClaimOfType("userInfo", userO)

	if err == nil {
		userInfo := new(UserInfo)
		database.DBConn.First(&userInfo, c.Params("id"))

		if userInfo.ID == 0 {
			return fiberUtils.SendJSONMessage("No User exists", false, 404)
		}

		if userO.User.Role == "Admin" || userO.ID == userInfo.ID {
			database.DBConn.Delete(&userInfo)
			return fiberUtils.SendSuccessResponse("User Successfully Deleted")
		}

		err = fiberUtils.SendJSONMessage("No permission to delete", false, 401)
	}

	return err
}

// GetUserInfoFromJWTClaim ...
func GetUserInfoFromJWTClaim(c *fiber.Ctx) UserInfo {
	fiberUtils.Ctx.New(c)
	userO := fiberUtils.GetJWTClaim("userInfo")
	userInfo := new(UserInfo)
	database.DBConn.First(&userInfo.ID, userO["id"])
	return *userInfo
}
