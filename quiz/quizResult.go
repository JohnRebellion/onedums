package quiz

import (
	"encoding/json"
	"fmt"
	"onedums/activity"
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
	err := fiberUtils.ParseBody(&quizresult)

	if err == nil {
		err = database.DBConn.Create(&quizresult).Error

		if err == nil {
			return fiberUtils.SendSuccessResponse("Created a new quiz result successfully")
		}
	}

	return err
}

// UpdateQuizResult ...
func UpdateQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizresult := new(QuizResult)
	err := fiberUtils.ParseBody(&quizresult)

	if err == nil {
		err = database.DBConn.Updates(&quizresult).Error

		if err == nil {
			return fiberUtils.SendSuccessResponse("Updated a quiz result successfully")
		}
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

// CheckQuizResult ...
func CheckQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizResult := new(QuizResult)
	err := fiberUtils.ParseBody(&quizResult)

	if err == nil {
		sum := 0
		quiz := new(Quiz)
		err = database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&quiz, quizResult.Quiz.ID).Error
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
								twilioService.SendSMS(fmt.Sprintf(".\n%s's grade on \"%s\": %2.f/100%%", student.UserInfo.User.Name, quiz.Title, quizResult.Percentage), student.Guardian.ContactNumber)
							}

							return fiberUtils.SendSuccessResponse("Created a new quiz result successfully")
						}
					}
				}
			}
		}
	}

	return err
}

// CheckUpdatedQuizResult ...
func CheckUpdatedQuizResult(c *fiber.Ctx) error {
	fiberUtils.Ctx.New(c)
	quizResult := new(QuizResult)
	err := fiberUtils.ParseBody(&quizResult)

	if err == nil {
		sum := 0
		quiz := new(Quiz)
		err = database.DBConn.Preload("Teacher.UserInfo.User").Preload("Subject").First(&quiz, quizResult.Quiz.ID).Error
		willSendSMS := false

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
						err = database.DBConn.Updates(&quizResult).Error

						if err == nil {
							if willSendSMS {
								twilioService.SendSMS(fmt.Sprintf(".\n%s's grade on \"%s\": %2.f/100%%", student.UserInfo.User.Name, quiz.Title, quizResult.Percentage), student.Guardian.ContactNumber)
							}

							return fiberUtils.SendSuccessResponse("Created a new quiz result successfully")
						}
					}
				}
			}
		}
	}

	return err
}

// GetQuizResultByStudentID ...
func GetQuizResultByStudentID(c *fiber.Ctx) error {
	quizResults := []QuizResult{}
	quizResultsFiltered := []QuizResult{}
	studentID, err := c.ParamsInt("studentId")

	if err == nil {
		err = database.DBConn.Preload("Quiz.Teacher.UserInfo.User").Preload("Quiz.Subject").Preload("Student.UserInfo.User").Preload("Student.Section").Find(&quizResults, "student_id = ?", studentID).Error
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
	}

	return err
}

// GetQuizResultByQuizID ...
func GetQuizResultByQuizID(c *fiber.Ctx) error {
	quizResults := []QuizResult{}
	quizResultsFiltered := []QuizResult{}
	quizID, err := c.ParamsInt("quizId")

	if err == nil {
		err = database.DBConn.Preload("Quiz.Teacher.UserInfo.User").Preload("Quiz.Subject").Preload("Student.UserInfo.User").Preload("Student.Section").Find(&quizResults, "quiz_id = ?", quizID).Error
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
	}

	return err
}

// GetStudentProgressBySubjectID ...
func GetStudentProgressBySubjectID(c *fiber.Ctx) error {
	studentID, err := c.ParamsInt("studentId")

	if err == nil {
		subjectID, err := c.ParamsInt("subjectId")

		if err == nil {
			quizzes := []Quiz{}
			err = database.DBConn.Preload("Subject").Find(&quizzes, "subject_id = ?", subjectID).Error

			if err == nil {
				quizResultsFiltered := []QuizResult{}
				activities := []activity.Activity{}

				for _, quiz := range quizzes {
					quizResults := []QuizResult{}
					err = database.DBConn.Preload("Quiz.Subject").Preload("Student.UserInfo.User").Find(&quizResults, "quiz_id = ? AND student_id = ?", quiz.ID, studentID).Error

					if err == nil {
						quizResultsFiltered = append(quizResultsFiltered, quizResults...)
					}
				}

				err = database.DBConn.Preload("Subject").Find(&activities, "subject_id = ?", subjectID).Error

				if err == nil {
					activityResultsFiltered := []activity.ActivityResult{}
					activityResults := []activity.ActivityResult{}

					for _, activity := range activities {
						err = database.DBConn.Preload("Activity.Subject").Preload("Student.UserInfo.User").Find(&activityResults, "activity_id = ? AND student_id = ?", activity.ID, studentID).Error

						if err == nil {
							activityResultsFiltered = append(activityResultsFiltered, activityResults...)
						}
					}

					return c.JSON(fiber.Map{
						"value": float64(len(quizResultsFiltered)+len(activityResultsFiltered)) / float64((len(quizzes))+len(activities)) * 100,
					})
				}
			}

		}
	}
	return err
}
