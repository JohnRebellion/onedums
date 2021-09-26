package learningMaterial

import (
	"fmt"
	"mime/multipart"
	"onedums/subject"
	"onedums/teacher"
	"onedums/user"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// LearningMaterial ...
type LearningMaterial struct {
	gorm.Model  `json:"-"`
	ID          uint            `json:"id" gorm:"primarykey"`
	Title       string          `json:"title"`
	Filename    string          `json:"filename"`
	Description string          `json:"description"`
	SubjectID   uint            `json:"-"`
	Subject     subject.Subject `json:"subject"`
	TeacherID   uint            `json:"-"`
	Teacher     teacher.Teacher `json:"teacher"`
}

// Item ...
type Item struct {
	Question         string   `json:"question"`
	Answer           string   `json:"answer"`
	Type             string   `json:"type"`
	IncorrectAnswers []string `json:"incorrectAnswers"`
}

// GetLearningMaterialsCurrentUser ...
func GetLearningMaterialsCurrentUser(c *fiber.Ctx) error {
	learningMaterials := []LearningMaterial{}
	learningMaterialsFiltered := []LearningMaterial{}
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").Find(&learningMaterials).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		for _, learningMaterial := range learningMaterials {
			_, err = os.Stat(learningMaterial.absoluteServerFilename())

			if learningMaterial.Teacher.UserInfo.User.ID != 0 &&
				learningMaterial.Subject.ID != 0 &&
				!os.IsNotExist(err) &&
				learningMaterial.Teacher.UserInfo.User.ID == userClaim.User.ID || userClaim.User.Role == "Admin" {
				learningMaterialsFiltered = append(learningMaterialsFiltered, learningMaterial)
			}
		}

		err = c.JSON(learningMaterialsFiltered)
	}

	return err
}

// GetLearningMaterialCurrentUser ...
func GetLearningMaterialCurrentUser(c *fiber.Ctx) error {
	learningMaterial := new(LearningMaterial)
	learningMaterialFiltered := new(LearningMaterial)
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&learningMaterial, c.Params("id")).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		_, err = os.Stat(learningMaterial.absoluteServerFilename())

		if learningMaterial.Teacher.UserInfo.User.ID != 0 &&
			learningMaterial.Subject.ID != 0 &&
			!os.IsNotExist(err) &&
			learningMaterial.Teacher.UserInfo.User.ID == userClaim.User.ID || userClaim.User.Role == "Admin" {
			learningMaterialFiltered = learningMaterial
		}

		err = c.JSON(learningMaterialFiltered)
	}

	return err
}

// UpdateLearningMaterialCurrentUser ...
func UpdateLearningMaterialCurrentUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	fiberUtils.ParseBody(learningMaterial)
	file, err := c.FormFile("filename")
	userClaim := user.GetUserInfoFromJWTClaim(c)
	teacher := new(teacher.Teacher)

	if err == nil {
		err = database.DBConn.Joins("UserInfo").Find(&teacher, "UserInfo.id = ?", userClaim.ID).Error

		if err == nil {
			if teacher.UserInfo.ID == userClaim.ID || userClaim.User.Role == "Admin" {
				learningMaterial.Filename = file.Filename
				c.SaveFile(file, learningMaterial.absoluteServerFilename())
				err := database.DBConn.Updates(&learningMaterial).Error

				if err == nil {
					return fiberUtils.SendSuccessResponse("Updated learning material successfully")
				}
			} else {
				return fiberUtils.SendJSONMessage("Learning material cannot be updated by current user", false, 401)
			}
		}
	}

	return err
}

// DeleteLearningMaterialCurrentUser ...
func DeleteLearningMaterialCurrentUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	userClaim := user.GetUserInfoFromJWTClaim(c)
	var err error
	teacher := new(teacher.Teacher)

	if err == nil {
		err = database.DBConn.Joins("UserInfo").Find(&teacher, "UserInfo.id = ?", userClaim.ID).Error

		if err == nil {
			if teacher.ID == learningMaterial.TeacherID || userClaim.User.Role == "Admin" {
				err = database.DBConn.Delete(learningMaterial, c.Params("id")).Error

				if err == nil {
					err = fiberUtils.SendSuccessResponse("Deleted learning material successfully")
				}
			} else {
				err = fiberUtils.SendJSONMessage("Learning material cannot be deleted by current user", false, 401)
			}
		}
	}

	return err
}

// UploadLearningMaterialCurrentUser ...
func UploadUpdatedLearningMaterialCurrentUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	file, err := c.FormFile("file")
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		id, err := strconv.ParseUint(c.FormValue("id"), 10, 32)
		learningMaterial.ID = uint(id)

		if err == nil {
			subjectID, err := strconv.ParseUint(c.FormValue("subjectId"), 10, 32)
			learningMaterial.Subject.ID = uint(subjectID)

			if err == nil {
				teacherID, err := strconv.ParseUint(c.FormValue("teacherId"), 10, 32)
				learningMaterial.Title = c.FormValue("title")
				learningMaterial.Description = c.FormValue("description")
				teacher := new(teacher.Teacher)

				if err == nil {
					err = database.DBConn.Joins("UserInfo").Find(&teacher, "UserInfo.id = ?", userClaim.ID).Error

					if err == nil {
						if userClaim.User.Role == "Admin" {
							learningMaterial.Teacher.ID = uint(teacherID)
						}

						learningMaterial.Filename = file.Filename
						err = uploadFile(c, file, learningMaterial)

						if err == nil {
							if teacher.ID == learningMaterial.TeacherID || userClaim.User.Role == "Admin" {
								err := database.DBConn.Updates(&learningMaterial).Error

								if err == nil {
									return fiberUtils.SendSuccessResponse("Uploaded updated learning material successfully")
								}
							} else {
								return fiberUtils.SendJSONMessage("Learning material cannot be uploaded by current user", false, 401)
							}
						}
					}
				}
			}
		}
	}

	return err
}

// UploadLearningMaterialCurrentUser ...
func UploadLearningMaterialCurrentUser(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	file, err := c.FormFile("file")
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		subjectID, err := strconv.ParseUint(c.FormValue("subjectId"), 10, 32)
		learningMaterial.Subject.ID = uint(subjectID)

		if err == nil {
			teacherID, err := strconv.ParseUint(c.FormValue("teacherId"), 10, 32)
			learningMaterial.Teacher.ID = uint(teacherID)
			learningMaterial.Title = c.FormValue("title")
			learningMaterial.Description = c.FormValue("description")
			teacher := new(teacher.Teacher)

			if err == nil {
				err = database.DBConn.Joins("UserInfo").Find(&teacher, "UserInfo.id = ?", userClaim.ID).Error

				if err == nil {
					if teacherID != uint64(teacher.ID) || userClaim.User.Role == "Admin" {
						teacherID = uint64(teacher.ID)
					}

					learningMaterial.Filename = file.Filename
					err = uploadFile(c, file, learningMaterial)

					if err == nil {
						if teacher.UserInfo.ID == userClaim.ID || userClaim.User.Role == "Admin" {
							err := database.DBConn.Updates(&learningMaterial).Error

							if err == nil {
								return fiberUtils.SendSuccessResponse("Uploaded updated learning material successfully")
							} else {
								return fiberUtils.SendJSONMessage("Learning material cannot be uploaded by current user", false, 401)
							}
						}
					}
				}
			}
		}
	}

	return err
}

// GetLearningMaterials ...
func GetLearningMaterials(c *fiber.Ctx) error {
	learningMaterials := []LearningMaterial{}
	learningMaterialsFiltered := []LearningMaterial{}
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").Find(&learningMaterials).Error

	if err == nil {
		for _, learningMaterial := range learningMaterials {
			_, err = os.Stat(learningMaterial.absoluteServerFilename())

			if learningMaterial.Teacher.UserInfo.User.ID != 0 &&
				learningMaterial.Subject.ID != 0 &&
				!os.IsNotExist(err) {
				learningMaterialsFiltered = append(learningMaterialsFiltered, learningMaterial)
			}
		}

		err = c.JSON(learningMaterialsFiltered)
	}

	return err
}

// GetLearningMaterial ...
func GetLearningMaterial(c *fiber.Ctx) error {
	learningMaterial := new(LearningMaterial)
	learningMaterialFiltered := new(LearningMaterial)
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&learningMaterial, c.Params("id")).Error

	if err == nil {
		_, err = os.Stat(learningMaterial.absoluteServerFilename())

		if learningMaterial.Teacher.UserInfo.User.ID != 0 &&
			learningMaterial.Subject.ID != 0 &&
			!os.IsNotExist(err) {
			learningMaterialFiltered = learningMaterial
		}

		err = c.JSON(learningMaterialFiltered)
	}

	return err
}

// NewLearningMaterial ...
func NewLearningMaterial(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	fiberUtils.ParseBody(&learningMaterial)

	err := database.DBConn.Create(&learningMaterial).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Created new learning material successfully")
	}

	return err
}

// UpdateLearningMaterial ...
func UpdateLearningMaterial(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	fiberUtils.ParseBody(learningMaterial)
	file, err := c.FormFile("filename")

	if err == nil {
		learningMaterial.Filename = file.Filename
		c.SaveFile(file, learningMaterial.absoluteServerFilename())
		err := database.DBConn.Updates(&learningMaterial).Error

		if err == nil {
			return fiberUtils.SendSuccessResponse("Updated learning material successfully")
		}
	}

	return err
}

// DeleteLearningMaterial ...
func DeleteLearningMaterial(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	err := database.DBConn.Delete(learningMaterial, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted learning material successfully")
	}

	return err
}

// UploadLearningMaterial ...
func UploadLearningMaterial(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	file, err := c.FormFile("file")

	if err == nil {
		subjectID, err := strconv.ParseUint(c.FormValue("subjectId"), 10, 32)
		learningMaterial.Subject.ID = uint(subjectID)

		if err == nil {
			teacherID, err := strconv.ParseUint(c.FormValue("teacherId"), 10, 32)
			learningMaterial.Teacher.ID = uint(teacherID)
			learningMaterial.Title = c.FormValue("title")
			learningMaterial.Description = c.FormValue("description")

			if err == nil {
				learningMaterial.Filename = file.Filename
				err = uploadFile(c, file, learningMaterial)

				if err == nil {
					err := database.DBConn.Create(&learningMaterial).Error

					if err == nil {
						return fiberUtils.SendSuccessResponse("Uploaded new learning material successfully")
					}
				}
			}
		}
	}

	return err
}

// UploadLearningMaterial ...
func UploadUpdatedLearningMaterial(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	learningMaterial := new(LearningMaterial)
	file, err := c.FormFile("file")

	if err == nil {
		id, err := strconv.ParseUint(c.FormValue("id"), 10, 32)
		learningMaterial.ID = uint(id)

		if err == nil {
			subjectID, err := strconv.ParseUint(c.FormValue("subjectId"), 10, 32)
			learningMaterial.Subject.ID = uint(subjectID)

			if err == nil {
				teacherID, err := strconv.ParseUint(c.FormValue("teacherId"), 10, 32)
				learningMaterial.Teacher.ID = uint(teacherID)
				learningMaterial.Title = c.FormValue("title")
				learningMaterial.Description = c.FormValue("description")

				if err == nil {
					learningMaterial.Filename = file.Filename
					err = uploadFile(c, file, learningMaterial)

					if err == nil {
						err := database.DBConn.Updates(&learningMaterial).Error

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

// DownloadLearningMaterialFile ...
func DownloadLearningMaterialFile(c *fiber.Ctx) error {
	learningMaterial := new(LearningMaterial)
	err := database.DBConn.Preload("Teacher.UserInfo").Preload("Subject").First(&learningMaterial, c.Params("id")).Error

	if err == nil {
		err = c.Download(learningMaterial.absoluteServerFilename())
	} else {
		return fiberUtils.SendJSONMessage("File not found", false, 404)
	}

	return err
}

func (learningMaterial *LearningMaterial) absoluteServerFilename() string {
	return fmt.Sprintf("files/learningMaterials/%d/%d/%s", learningMaterial.Teacher.ID, learningMaterial.Subject.ID, learningMaterial.Filename)
}

func (learningMaterial *LearningMaterial) absoluteServerDirectory() string {
	return fmt.Sprintf("files/learningMaterials/%d/%d", learningMaterial.Teacher.ID, learningMaterial.Subject.ID)
}

func uploadFile(c *fiber.Ctx, file *multipart.FileHeader, learningMaterial *LearningMaterial) error {
	err := makeDirectoryIfNotExists(learningMaterial.absoluteServerDirectory())

	if err == nil {
		err = c.SaveFile(file, learningMaterial.absoluteServerFilename())
	}

	return err
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModeDir|0755)
	}

	return nil
}
