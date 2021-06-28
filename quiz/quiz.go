// Package quiz provides quiz CRUD
package quiz

import (
	"onedums/subject"
	"onedums/teacher"
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
	Teacher          teacher.Teacher `json:"teacher"`
	SubjectID        uint            `json:"-"`
	Subject          subject.Subject `json:"subject"`
	Title            string          `json:"title"`
	Items            datatypes.JSON  `json:"items"`
	DateOfSubmission time.Time       `json:"dateOfSubmission"`
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
	database.DBConn.Preload("Teacher.User").Preload("Subject").Find(&quizzes)
	return c.JSON(quizzes)
}

// GetQuiz ...
func GetQuiz(c *fiber.Ctx) error {
	quiz := new(Quiz)
	database.DBConn.Preload("Teacher.User").Preload("Subject").First(&quiz, c.Params("id"))
	return c.JSON(quiz)
}

// NewQuiz ...
func NewQuiz(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quiz := new(Quiz)
	fiberUtils.ParseBody(&quiz)
	database.DBConn.Create(&quiz)
	return fiberUtils.SendSuccessResponse("Created new quiz successfully")
}

// UpdateQuiz ...
func UpdateQuiz(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quiz := new(Quiz)
	fiberUtils.ParseBody(quiz)
	database.DBConn.Updates(&quiz)
	return fiberUtils.SendSuccessResponse("Updated quiz successfully")
}

// DeleteQuiz ...
func DeleteQuiz(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	database.DBConn.Delete(&Quiz{}, c.Params("id"))
	return fiberUtils.SendSuccessResponse("Deleted quiz successfully")
}
