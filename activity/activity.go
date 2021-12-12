package activity

import (
	"fmt"
	"mime/multipart"
	"onedums/subject"
	"onedums/teacher"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// Activity ...
type Activity struct {
	gorm.Model       `json:"-"`
	ID               uint            `json:"id" gorm:"primarykey"`
	TeacherID        uint            `json:"-"`
	Teacher          teacher.Teacher `json:"teacher" gorm:"constraint:OnDelete:SET NULL;"`
	SubjectID        uint            `json:"-"`
	Subject          subject.Subject `json:"subject" gorm:"constraint:OnDelete:SET NULL;"`
	Title            string          `json:"title"`
	Filename         string          `json:"filename"`
	DateOfSubmission time.Time       `json:"dateOfSubmission"`
	CreatedAt        time.Time       `json:"dateCreated"`
}

// GetActivities ...
func GetActivities(c *fiber.Ctx) error {
	activities := []Activity{}
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").Find(&activities).Error

	if err == nil {
		err = c.JSON(activities)
	}

	return err
}

// GetActivity ...
func GetActivity(c *fiber.Ctx) error {
	activity := new(Activity)
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&activity, c.Params("id")).Error

	if err == nil {
		err = c.JSON(&activity)
	}

	return err
}

// NewActivity ...
func NewActivity(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activity := new(Activity)
	err := fiberUtils.ParseBody(&activity)

	if err == nil {
		err = database.DBConn.Create(&activity).Error

		if err == nil {
			return fiberUtils.SendSuccessResponse("Created a new activity successfully")
		}
	}

	return err
}

// UpdateActivity ...
func UpdateActivity(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activity := new(Activity)
	err := fiberUtils.ParseBody(&activity)

	if err == nil {
		err = database.DBConn.Updates(&activity).Error

		if err == nil {
			return fiberUtils.SendSuccessResponse("Updated a activity successfully")
		}
	}

	return err
}

// DeleteActivity ...
func DeleteActivity(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&Activity{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted a activity successfully")
	}

	return err
}

// UploadActivity ...
func UploadActivity(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activity := new(Activity)
	file, err := c.FormFile("file")

	if err == nil {
		subjectID, err := strconv.ParseUint(c.FormValue("subjectId"), 10, 32)
		activity.Subject.ID = uint(subjectID)

		if err == nil {
			teacherID, err := strconv.ParseUint(c.FormValue("teacherId"), 10, 32)
			activity.Teacher.ID = uint(teacherID)
			activity.Title = c.FormValue("title")
			dateOfSubmission, err := time.Parse("2006-01-02T15:04:05.999Z07:00", c.FormValue("dateOfSubmission"))

			if err == nil {
				activity.DateOfSubmission = dateOfSubmission

				if err == nil {
					activity.Filename = file.Filename
					err = uploadFile(c, file, activity)

					if err == nil {
						err := database.DBConn.Create(&activity).Error

						if err == nil {
							return fiberUtils.SendSuccessResponse("Uploaded new learning material successfully")
						}
					}
				}
			}
		}
	}

	return err
}

// UploadActivity ...
func UploadUpdatedActivity(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	activity := new(Activity)
	file, err := c.FormFile("file")

	if err == nil {
		id, err := strconv.ParseUint(c.FormValue("id"), 10, 32)
		activity.ID = uint(id)

		if err == nil {
			subjectID, err := strconv.ParseUint(c.FormValue("subjectId"), 10, 32)
			activity.Subject.ID = uint(subjectID)

			if err == nil {
				teacherID, err := strconv.ParseUint(c.FormValue("teacherId"), 10, 32)
				activity.Teacher.ID = uint(teacherID)
				activity.Title = c.FormValue("title")
				dateOfSubmission, err := time.Parse("2006-01-02T15:04:05.999Z07:00", c.FormValue("dateOfSubmission"))

				if err == nil {
					activity.DateOfSubmission = dateOfSubmission

					if err == nil {
						activity.Filename = file.Filename
						err = uploadFile(c, file, activity)

						if err == nil {
							err := database.DBConn.Updates(&activity).Error

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

// DownloadActivityFile ...
func DownloadActivityFile(c *fiber.Ctx) error {
	activity := new(Activity)
	err := database.DBConn.Preload("Teacher.UserInfo").Preload("Subject").First(&activity, c.Params("id")).Error

	if err == nil {
		err = c.Download(activity.absoluteServerFilename())
	} else {
		return fiberUtils.SendJSONMessage("File not found", false, 404)
	}

	return err
}

func (activity *Activity) absoluteServerFilename() string {
	return fmt.Sprintf("files/activities/%d/%d/%s", activity.Teacher.ID, activity.Subject.ID, activity.Filename)
}

func (activity *Activity) absoluteServerDirectory() string {
	return fmt.Sprintf("files/activities/%d/%d", activity.Teacher.ID, activity.Subject.ID)
}

func uploadFile(c *fiber.Ctx, file *multipart.FileHeader, activity *Activity) error {
	err := makeDirectoryIfNotExists(activity.absoluteServerDirectory())

	if err == nil {
		err = c.SaveFile(file, activity.absoluteServerFilename())
	}

	return err
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModeDir|0755)
	}

	return nil
}

// GetActivitiesBySubjectID ...
func GetActivitiesBySubjectID(c *fiber.Ctx) error {
	activities := []Activity{}
	subjectID, err := c.ParamsInt("subjectId")

	if err == nil {
		err = database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").Find(&activities, "subject_id = ?", subjectID).Error

		if err == nil {
			return c.JSON(activities)
		}
	}

	return err
}
