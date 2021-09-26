// Package quiz provides quiz CRUD
package quiz

import (
	"onedums/subject"
	"onedums/teacher"
	"onedums/user"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// Quiz ...
type Quiz struct {
	gorm.Model       `json:"-"`
	ID               uint            `json:"id" gorm:"primarykey"`
	TeacherID        uint            `json:"-"`
	Teacher          teacher.Teacher `json:"teacher" gorm:"constraint:OnDelete:SET NULL;"`
	SubjectID        uint            `json:"-"`
	Subject          subject.Subject `json:"subject" gorm:"constraint:OnDelete:SET NULL;"`
	Title            string          `json:"title"`
	Items            datatypes.JSON  `json:"items"`
	DateOfSubmission time.Time       `json:"dateOfSubmission"`
	TimeLimit        uint64          `json:"timeLimit"`
}

// Item ...
type Item struct {
	Question         string   `json:"question"`
	Answer           string   `json:"answer"`
	Type             string   `json:"type"`
	IncorrectAnswers []string `json:"incorrectAnswers"`
}

// GetQuizzes ...
func GetQuizzes(c *fiber.Ctx) error {
	quizzes := []Quiz{}
	quizzesFiltered := []Quiz{}
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").Find(&quizzes).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		for _, quiz := range quizzes {
			if quiz.Teacher.UserInfo.User.ID != 0 &&
				quiz.Subject.ID != 0 ||
				userClaim.User.Role == "Admin" {
				quizzesFiltered = append(quizzesFiltered, quiz)
			}
		}

		err = c.JSON(quizzesFiltered)
	}

	return err
}

// GetQuiz ...
func GetQuiz(c *fiber.Ctx) error {
	quiz := new(Quiz)
	quizFiltered := new(Quiz)
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&quiz, c.Params("id")).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		if quiz.Teacher.UserInfo.User.ID != 0 &&
			quiz.Subject.ID != 0 ||
			userClaim.User.Role == "Admin" {
			quizFiltered = quiz
		}

		err = c.JSON(quizFiltered)
	}

	return err
}

// NewQuiz ...
func NewQuiz(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quiz := new(Quiz)
	fiberUtils.ParseBody(&quiz)
	err := database.DBConn.Create(&quiz).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Created new quiz successfully")
	}

	return err
}

// UpdateQuiz ...
func UpdateQuiz(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quiz := new(Quiz)
	fiberUtils.ParseBody(quiz)
	err := database.DBConn.Updates(&quiz).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Updated quiz successfully")
	}

	return err
}

// DeleteQuiz ...
func DeleteQuiz(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&Quiz{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted quiz successfully")
	}

	return err
}
