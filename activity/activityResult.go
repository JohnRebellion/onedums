package activity

import (
	"fmt"
	"log"
	"mime/multipart"
	"onedums/student"
	"onedums/twilioService"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// ActivityResult ...
type ActivityResult struct {
	gorm.Model  `json:"-"`
	ID          uint            `json:"id" gorm:"primarykey"`
	ActivityID  uint            `json:"-"`
	Activity    Activity        `json:"activity" gorm:"constraint:OnDelete:SET NULL;"`
	StudentID   uint            `json:"-"`
	Student     student.Student `json:"student" gorm:"constraint:OnDelete:SET NULL;"`
	Comment     string          `json:"comment"`
	IsSubmitted bool            `json:"isSubmitted"`
	Percentage  float64         `json:"percentage"`
	Filename    string          `json:"filename"`
}

// GetActivityResults ...
func GetActivityResults(c *fiber.Ctx) error {
	activityResults := []ActivityResult{}
	err := database.DBConn.Preload("Student.UserInfo.User").Preload("Activity").Find(&activityResults).Error

	if err == nil {
		err = c.JSON(activityResults)
	}

	return err
}

// GetActivityResult ...
func GetActivityResult(c *fiber.Ctx) error {
	activityResult := new(ActivityResult)
	err := database.DBConn.Preload("Student.UserInfo.User").Preload("Activity").First(&activityResult, c.Params("id")).Error

	if err == nil {
		err = c.JSON(&activityResult)
	}

	return err
}

// NewActivityResult ...
func NewActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activityResult := new(ActivityResult)
	err := fiberUtils.ParseBody(&activityResult)

	if err == nil {
		err = database.DBConn.Create(&activityResult).Error

		if err == nil {
			return fiberUtils.SendSuccessResponse("Created a new activityResult successfully")
		}
	}

	return err
}

// UpdateActivityResult ...
func UpdateActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activityResult := new(ActivityResult)
	err := fiberUtils.ParseBody(&activityResult)

	if err == nil {
		err := database.DBConn.Updates(&activityResult).Error

		if err == nil {
			return fiberUtils.SendSuccessResponse("Updated a activityResult successfully")
		}
	}

	return err
}

// DeleteActivityResult ...
func DeleteActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&ActivityResult{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted a activityResult successfully")
	}

	return err
}

// UploadActivityResult ...
func UploadActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activityResult := new(ActivityResult)
	file, err := c.FormFile("file")

	if err == nil {
		studentID, err := strconv.ParseUint(c.FormValue("studentId"), 10, 32)
		activityResult.Student.ID = uint(studentID)

		if err == nil {
			activityID, err := strconv.ParseUint(c.FormValue("activityId"), 10, 32)
			activityResult.Activity.ID = uint(activityID)

			if err == nil {
				activityResult.Filename = file.Filename
				err = uploadResultFile(c, file, activityResult)

				if err == nil {
					err := database.DBConn.Create(&activityResult).Error

					if err == nil {
						return fiberUtils.SendSuccessResponse("Uploaded new learning material successfully")
					}
				}
			}
		}
	}

	return err
}

// UploadActivityResult ...
func UploadUpdatedActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activityResult := new(ActivityResult)
	file, err := c.FormFile("file")

	if err == nil {
		activityResult.Comment = c.FormValue("comment")
		activityResult.Percentage, err = strconv.ParseFloat(c.FormValue("percentage"), 64)

		if err == nil {
			id, err := strconv.ParseUint(c.FormValue("id"), 10, 32)
			activityResult.ID = uint(id)

			if err == nil {
				studentID, err := strconv.ParseUint(c.FormValue("studentId"), 10, 32)
				activityResult.Student.ID = uint(studentID)

				if err == nil {
					activityID, err := strconv.ParseUint(c.FormValue("activityId"), 10, 32)
					activityResult.Activity.ID = uint(activityID)

					if err == nil {
						activityResult.Filename = file.Filename
						err = uploadResultFile(c, file, activityResult)

						if err == nil {
							err := database.DBConn.Updates(&activityResult).Error

							if err == nil {
								return fiberUtils.SendSuccessResponse("Uploaded updated learning material successfully")
							}
						}
					}
				}
			}
		}
	}

	return err
}

// DownloadActivityResultFile ...
func DownloadActivityResultFile(c *fiber.Ctx) error {
	activityResult := new(ActivityResult)
	err := database.DBConn.Preload("Student.UserInfo").Preload("Activity").First(&activityResult, c.Params("id")).Error

	if err == nil {
		err = c.Download(activityResult.absoluteServerFilename())
	} else {
		return fiberUtils.SendJSONMessage("File not found", false, 404)
	}

	return err
}

func (activityResult *ActivityResult) absoluteServerFilename() string {
	return fmt.Sprintf("files/activity/%d/results/%d/%s", activityResult.Activity.ID, activityResult.Student.ID, activityResult.Filename)
}

func (activityResult *ActivityResult) absoluteServerDirectory() string {
	return fmt.Sprintf("files/activity/%d/results/%d", activityResult.Activity.ID, activityResult.Student.ID)
}

// GetActivityResultsByActivityID ...
func GetActivityResultsByActivityID(c *fiber.Ctx) error {
	activityResults := []ActivityResult{}
	activityID, err := c.ParamsInt("activityId")

	if err == nil {
		err = database.DBConn.Preload("Student.UserInfo.User").Preload("Activity").Find(&activityResults, "activity_id = ?", activityID).Error

		if err == nil {
			return c.JSON(activityResults)
		}
	}

	return err
}

// CheckActivityResult ...
func CheckActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activityResult := new(ActivityResult)
	err := fiberUtils.ParseBody(&activityResult)

	if err == nil {
		id := activityResult.ID
		comment := activityResult.Comment
		percentage := activityResult.Percentage
		willSendSMS := true
		err = database.DBConn.Preload("Student.UserInfo.User").Preload("Activity").First(&activityResult, id).Error

		if err == nil {
			activityResult.Comment = comment
			activityResult.Percentage = percentage
			err = database.DBConn.Updates(&activityResult).Error

			if err == nil {
				if willSendSMS {
					twilioService.SendSMS(fmt.Sprintf(".\n%s's grade on \"%s\": %2.f%s", activityResult.Student.UserInfo.User.Name, activityResult.Activity.Title, activityResult.Percentage, "%"), activityResult.Student.Guardian.ContactNumber)
				} else {
					log.Printf(".\n%s's grade on \"%s\": %2.f%s", activityResult.Student.UserInfo.User.Name, activityResult.Activity.Title, activityResult.Percentage, "%")
				}

				return fiberUtils.SendSuccessResponse("Created a new activity result successfully")
			}
		}
	}

	return err
}

func uploadResultFile(c *fiber.Ctx, file *multipart.FileHeader, activityResult *ActivityResult) error {
	err := makeDirectoryIfNotExists(activityResult.absoluteServerDirectory())

	if err == nil {
		err = c.SaveFile(file, activityResult.absoluteServerFilename())
	}

	return err
}
