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
	User          User      `json:"user" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Birthday      time.Time `json:"birthday"`
	Gender        string    `json:"gender"`
	Address       string    `json:"address"`
	ContactNumber string    `json:"contactNumber"`
}

// GetUserInfo Get a User by id
func GetUserInfo(c *fiber.Ctx) error {
	userInfo := new(UserInfo)
	userInfoFiltered := new(UserInfo)
	err := database.DBConn.Preload("User").Find(&userInfo, c.Params("id")).Error
	userClaim := GetUserInfoFromJWTClaim(c)

	if err == nil {
		if userInfo.User.ID != 0 ||
			userClaim.User.Role == "Admin" {
			userInfoFiltered = userInfo
		}

		err = c.JSON(userInfoFiltered)
	}

	return err
}

// GetUserInfos Get all users
func GetUserInfos(c *fiber.Ctx) error {
	userInfos := []UserInfo{}
	userInfosFiltered := []UserInfo{}
	err := database.DBConn.Preload("User").Find(&userInfos).Error
	userClaim := GetUserInfoFromJWTClaim(c)

	if err == nil {
		for _, userInfo := range userInfos {
			if userInfo.User.ID != 0 ||
				userClaim.User.Role == "Admin" {
				userInfosFiltered = append(userInfosFiltered, userInfo)
			}
		}

		err = c.JSON(userInfosFiltered)
	}

	return err
}

// NewUserInfo Creating a User
func NewUserInfo(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userInfo := new(UserInfo)
	err := fiberUtils.ParseBody(userInfo)

	if err == nil {
		var usernameCount int64

		if database.DBConn.Model(&userInfo.User).Where("username = ?", userInfo.User.Username).Count(&usernameCount).Error == nil {
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
				if database.DBConn.Create(&userInfo).Error == nil {
					return fiberUtils.SendJSONMessage("User Successfully Created", true, 201)
				}
			}
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

		userInfos := []UserInfo{}

		if database.DBConn.Preload("User").Find(&userInfos).Error == nil {
			for _, userInfo := range userInfos {
				if userInfo.User.Username == user.Username {
					if passwordHashing.CheckPasswordHash(password, userInfo.User.Password) {
						_, err := fiberUtils.GenerateJWTSignedString(fiber.Map{"userInfo": userInfo, "isAdmin": userInfo.User.Role == "Admin"})

						if err == nil {
							return fiberUtils.SendSuccessResponse("Access Granted!")
						}
					} else {
						return fiberUtils.SendJSONMessage("Incorrect Username or Password", false, 401)
					}
				}
			}
		}
	}

	return err
}

// UpdateUserInfo Update User
func UpdateUserInfo(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userClaim := GetUserInfoFromJWTClaim(c)
	userInfo := new(UserInfo)
	err := fiberUtils.ParseBody(userInfo)
	userInfo.User.ID = userInfo.ID

	if err == nil {
		if userClaim.User.Role == "Admin" || userClaim.ID == userInfo.ID {
			if len(userInfo.User.Username) == 0 || len(userInfo.User.Password) == 0 || len(userInfo.User.Name) == 0 {
				return fiberUtils.SendBadRequestResponse("Please Input Username, Password and Name")
			}

			if len(userInfo.User.Username) < 3 || len(userInfo.User.Password) < 8 || len(userInfo.User.Name) < 3 {
				return fiberUtils.SendBadRequestResponse("Required Mininum Length of Username, Name and Password is 3, 3 and 8 respectively")
			}

			userInfo.User.Password, err = passwordHashing.HashPassword(userInfo.User.Password)
			fiberUtils.LogError(err)
			var existingUserInfo UserInfo

			if database.DBConn.First(&existingUserInfo, userInfo.ID).Error == nil {
				if existingUserInfo.ID == 0 {
					return fiberUtils.SendJSONMessage("No User exists", false, 404)
				}

				if database.DBConn.Updates(&userInfo).Error == nil {
					return fiberUtils.SendSuccessResponse("User Successfully Updated")
				}
			}
		}

		return fiberUtils.SendJSONMessage("No permission to update", false, 401)
	}

	return err
}

// DeleteUserInfo Delete User by id
func DeleteUserInfo(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	userClaim := GetUserInfoFromJWTClaim(c)
	userInfo := new(UserInfo)
	err := database.DBConn.First(&userInfo, c.Params("id")).Error

	if err == nil {
		if userInfo.ID == 0 {
			return fiberUtils.SendJSONMessage("No User exists", false, 404)
		}

		if userClaim.User.Role == "Admin" || userClaim.ID == userInfo.ID {
			if database.DBConn.Delete(&userInfo).Error == nil {
				return fiberUtils.SendSuccessResponse("User Successfully Deleted")
			}
		}

		return fiberUtils.SendJSONMessage("No permission to delete", false, 401)
	}

	return err
}

// GetUserInfoFromJWTClaim ...
func GetUserInfoFromJWTClaim(c *fiber.Ctx) UserInfo {
	fiberUtils.Ctx.New(c)
	userO := fiberUtils.GetJWTClaim("userInfo")
	userInfo := new(UserInfo)
	err := database.DBConn.Preload("User").First(&userInfo, userO["id"]).Error
	if err == nil {
	}

	return *userInfo
}
