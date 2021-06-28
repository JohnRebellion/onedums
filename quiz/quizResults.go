package quiz

import (
	"encoding/json"
	"onedums/student"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
)

// QuizResult ...
type QuizResult struct {
	gorm.Model  `json:"-"`
	ID          uint            `json:"id" gorm:"primarykey"`
	QuizID      uint            `json:"-"`
	Quiz        Quiz            `json:"quiz"`
	StudentID   uint            `json:"-"`
	Student     student.Student `json:"student"`
	Answers     datatypes.JSON  `json:"answers"`
	IsSubmitted bool            `json:"isSubmitted"`
	Percentage  float64         `json:"percentage"`
}

// GetQuizResults ...
func GetQuizResults(c *fiber.Ctx) error {
	quizresults := []QuizResult{}
	database.DBConn.Preload("Quiz.Teacher.User").Preload("Quiz.Subject").Preload("Student.User").Preload("Student.Section").Find(&quizresults)
	return c.JSON(quizresults)
}

// GetQuizResult ...
func GetQuizResult(c *fiber.Ctx) error {
	quizresult := new(QuizResult)
	database.DBConn.Preload("Quiz.Teacher.User").Preload("Quiz.Subject").Preload("Student.User").Preload("Student.Section").First(&quizresult, c.Params("id"))
	return c.JSON(&quizresult)
}

// NewQuizResult ...
func NewQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizresult := new(QuizResult)
	fiberUtils.ParseBody(&quizresult)
	database.DBConn.Create(&quizresult)
	return fiberUtils.SendSuccessResponse("Created a new quiz result successfully")
}

// UpdateQuizResult ...
func UpdateQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizresult := new(QuizResult)
	fiberUtils.ParseBody(&quizresult)
	database.DBConn.Updates(&quizresult)
	return fiberUtils.SendSuccessResponse("Updated a quiz result successfully")
}

// DeleteQuizResult ...
func DeleteQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	database.DBConn.Delete(&QuizResult{}, c.Params("id"))
	return fiberUtils.SendSuccessResponse("Deleted a quiz result successfully")
}

// NewQuizResult ...
func CheckQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizResult := new(QuizResult)
	fiberUtils.ParseBody(&quizResult)
	sum := 0
	quiz := new(Quiz)
	database.DBConn.First(&quiz, quizResult.Quiz.ID)
	items := []Item{}
	err := json.Unmarshal(quiz.Items, &items)

	if err == nil {
		answers := []string{}
		err = json.Unmarshal(quizResult.Answers, &answers)

		if err == nil {
			for i, item := range items {
				if item.Answer == answers[i] {
					sum++
				}
			}

			quizResult.Percentage = float64(sum) * 100 / float64(len(items))
			database.DBConn.Create(&quizResult)
			err = fiberUtils.SendSuccessResponse("Created a new quiz result successfully")
		}
	}

	return err
}
