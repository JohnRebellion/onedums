package quiz

import (
	"encoding/json"
	"fmt"
	"onedums/student"
	"onedums/twilioService"
	"onedums/user"

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
	Quiz        Quiz            `json:"quiz" gorm:"constraint:OnDelete:SET NULL;"`
	StudentID   uint            `json:"-"`
	Student     student.Student `json:"student" gorm:"constraint:OnDelete:SET NULL;"`
	Answers     datatypes.JSON  `json:"answers"`
	Comment     string          `json:"comment"`
	IsSubmitted bool            `json:"isSubmitted"`
	Percentage  float64         `json:"percentage"`
}

// GetQuizResults ...
func GetQuizResults(c *fiber.Ctx) error {
	quizResults := []QuizResult{}
	quizResultsFiltered := []QuizResult{}
	err := database.DBConn.Preload("Quiz.Teacher.UserInfo.User").Preload("Quiz.Subject").Preload("Student.UserInfo.User").Preload("Student.Section").Find(&quizResults).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		for _, quizResult := range quizResults {
			if quizResult.Quiz.Teacher.UserInfo.User.ID != 0 &&
				quizResult.Quiz.Subject.ID != 0 &&
				quizResult.Student.UserInfo.User.ID != 0 &&
				quizResult.Student.Section.ID != 0 ||
				userClaim.User.Role == "Admin" {
				quizResultsFiltered = append(quizResultsFiltered, quizResult)
			}
		}

		err = c.JSON(quizResultsFiltered)
	}

	return err
}

// GetQuizResult ...
func GetQuizResult(c *fiber.Ctx) error {
	quizResult := new(QuizResult)
	quizResultFiltered := new(QuizResult)
	err := database.DBConn.Preload("Quiz.Teacher.UserInfo.User").Preload("Quiz.Subject").Preload("Student.UserInfo.User").Preload("Student.Section").First(&quizResult, c.Params("id")).Error
	userClaim := user.GetUserInfoFromJWTClaim(c)

	if err == nil {
		if quizResult.Quiz.Teacher.UserInfo.User.ID != 0 &&
			quizResult.Quiz.Subject.ID != 0 &&
			quizResult.Student.UserInfo.User.ID != 0 &&
			quizResult.Student.Section.ID != 0 ||
			userClaim.User.Role == "Role" {
			quizResultFiltered = quizResult
		}

		err = c.JSON(&quizResultFiltered)
	}

	return err
}

// NewQuizResult ...
func NewQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizresult := new(QuizResult)
	fiberUtils.ParseBody(&quizresult)
	err := database.DBConn.Create(&quizresult).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Created a new quiz result successfully")
	}

	return err
}

// UpdateQuizResult ...
func UpdateQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizresult := new(QuizResult)
	fiberUtils.ParseBody(&quizresult)
	err := database.DBConn.Updates(&quizresult).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Updated a quiz result successfully")
	}

	return err
}

// DeleteQuizResult ...
func DeleteQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	err := database.DBConn.Delete(&QuizResult{}, c.Params("id")).Error

	if err == nil {
		return fiberUtils.SendSuccessResponse("Deleted a quiz result successfully")
	}

	return err
}

// NewQuizResult ...
func CheckQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizResult := new(QuizResult)
	fiberUtils.ParseBody(&quizResult)
	sum := 0
	quiz := new(Quiz)
	err := database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&quiz, quizResult.Quiz.ID).Error
	willSendSMS := true

	if err == nil {
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
				student := new(student.Student)
				err = database.DBConn.Preload("Section").Preload("UserInfo.User").Find(&student, quizResult.Student.ID).Error

				if err == nil {
					err = database.DBConn.Create(&quizResult).Error

					if err == nil {
						if willSendSMS {
							twilioService.SendSMS(fmt.Sprintf(".\n%s's grade on \"%s\": %2.f%s", student.UserInfo.User.Name, quiz.Title, quizResult.Percentage, "%"), quizResult.Student.Guardian.ContactNumber)
						}

						return fiberUtils.SendSuccessResponse("Created a new quiz result successfully")
					}
				}
			}
		}
	}

	return err
}
