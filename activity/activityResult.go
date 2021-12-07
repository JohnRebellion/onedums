package activity

import (
	"fmt"
	"mime/multipart"
	"onedums/student"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
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
	Answers     datatypes.JSON  `json:"answers"`
	Comment     string          `json:"comment"`
	IsSubmitted bool            `json:"isSubmitted"`
	Percentage  float64         `json:"percentage"`
	Filename    string          `json:"filename"`
}

// GetActivityResults ...
func GetActivityResults(c *fiber.Ctx) error {
	activityResults := []ActivityResult{}
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").Find(&activityResults).Error

	if err == nil {
		err = c.JSON(activityResults)
	}

	return err
}

// GetActivityResult ...
func GetActivityResult(c *fiber.Ctx) error {
	activityResult := new(ActivityResult)
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&activityResult, c.Params("id")).Error

	if err == nil {
		err = c.JSON(&activityResult)
	}

	return err
}

// NewActivityResult ...
func NewActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activityResult := new(ActivityResult)
	fiberUtils.ParseBody(&activityResult)
	err := database.DBConn.Create(&activityResult).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Created a new activityResult successfully")
	}

	return err
}

// UpdateActivityResult ...
func UpdateActivityResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activityResult := new(ActivityResult)
	fiberUtils.ParseBody(&activityResult)
	err := database.DBConn.Updates(&activityResult).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Updated a activityResult successfully")
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

	return err
}

// DownloadActivityResultFile ...
func DownloadActivityResultFile(c *fiber.Ctx) error {
	activityResult := new(ActivityResult)
	err := database.DBConn.Preload("Teacher.UserInfo").Preload("Subject").First(&activityResult, c.Params("id")).Error

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

// GetActivityResultsBySubjectID ...
func GetActivityResultsBySubjectID(c *fiber.Ctx) error {
	activityResults := []ActivityResult{}
	subjectID, err := c.ParamsInt("subjectId")

	if err == nil {
		err = database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").Find(&activityResults, "subject_id = ?", subjectID).Error

		if err == nil {
			return c.JSON(activityResults)
		}
	}

	return err
}

// CheckActivityResult ...
// func CheckActivityResult(c *fiber.Ctx) error {
// 	fiberUtils.Ctx.New(c)
// 	activityResult := new(ActivityResult)
// 	fiberUtils.ParseBody(&activityResult)
// 	sum := 0
// 	activity := new(Activity)
// 	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&activity, activityResult.Activity.ID).Error
// 	willSendSMS := true

// 	if err == nil {
// 		items := []Item{}
// 		err := json.Unmarshal(activity.Items, &items)

// 		if err == nil {
// 			answers := []string{}
// 			err = json.Unmarshal(activityResult.Answers, &answers)

// 			if err == nil {
// 				for i, item := range items {
// 					if item.Answer == answers[i] {
// 						sum++
// 					}
// 				}

// 				activityResult.Percentage = float64(sum) * 100 / float64(len(items))
// 				student := new(student.Student)
// 				err = database.DBConn.Preload("Section").Preload("UserInfo.User").Find(&student, activityResult.Student.ID).Error

// 				if err == nil {
// 					err = database.DBConn.Create(&activityResult).Error

// 					if err == nil {
// 						if willSendSMS {
// 							twilioService.SendSMS(fmt.Sprintf(".\n%s's grade on \"%s\": %2.f%s", student.UserInfo.User.Name, activity.Title, activityResult.Percentage, "%"), student.Guardian.ContactNumber)
// 						}

// 						return fiberUtils.SendSuccessResponse("Created a new activity result successfully")
// 					}
// 				}
// 			}
// 		}
// 	}

// 	return err
// }

func uploadResultFile(c *fiber.Ctx, file *multipart.FileHeader, activityResult *ActivityResult) error {
	err := makeDirectoryIfNotExists(activityResult.absoluteServerDirectory())

	if err == nil {
		err = c.SaveFile(file, activityResult.absoluteServerFilename())
	}

	return err
}
