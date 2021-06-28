package main

import (
	"fmt"
	"log"
	"onedums/envRouting"
	"onedums/quiz"
	"onedums/section"
	"onedums/student"
	"onedums/subject"
	"onedums/teacher"
	"onedums/user"
	"time"

	"github.com/JohnRebellion/go-utils/database"
	fiberUtils "github.com/JohnRebellion/go-utils/fiber"
	"github.com/JohnRebellion/go-utils/passwordHashing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	envRouting.LoadEnv()
	// database.SQLiteConnect(envRouting.SQLiteFilename)
	database.MySQLConnect(envRouting.MySQLUsername, envRouting.MySQLPassword, envRouting.MySQLHost, envRouting.DatabaseName)
	err := database.DBConn.AutoMigrate(
		&section.Section{},
		&subject.Subject{},
		&student.Student{},
		&teacher.Teacher{},
		&user.User{},
		&user.UserInfo{},
		&quiz.Quiz{},
		&quiz.QuizResult{},
	)

	var existingUserInfo user.UserInfo
	database.DBConn.First(&existingUserInfo, 1)
	password, err := passwordHashing.HashPassword("12345678")

	if existingUserInfo.ID == 0 {
		database.DBConn.Create(&user.UserInfo{
			User: user.User{
				Username: "admin",
				Password: password,
				Name:     "Admin",
				Role:     "Admin",
			},
		})
	}

	if err != nil {
		log.Fatal(err.Error())
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	endpoint := app.Group("/api/v1/", logger.New())
	userInfoEndpoint := endpoint.Group("/userInfo", logger.New())
	userInfoEndpoint.Post("/", user.NewUserInfo)
	userInfoEndpoint.Post("/auth", user.Authenticate)

	app.Static("/", envRouting.StaticWebLocation)
	app.Use(fiberUtils.AuthenticationMiddleware(fiberUtils.JWTConfig{
		Duration:     15 * time.Minute,
		CookieMaxAge: 15 * 60,
		SetCookies:   true,
		SecretKey:    []byte(envRouting.SecretKey),
	}))

	userInfoEndpoint.Get("/:id", user.GetUserInfo)
	userInfoEndpoint.Get("/", user.GetUserInfos)
	userInfoEndpoint.Put("/", user.UpdateUserInfo)
	userInfoEndpoint.Delete("/:id", user.DeleteUserInfo)

	userEndpoint := endpoint.Group("/user", logger.New())
	userEndpoint.Get("/", user.GetUsers)
	userEndpoint.Get("/:id", user.GetUser)
	userEndpoint.Post("/", user.NewUser)
	userEndpoint.Put("/", user.UpdateUser)
	userEndpoint.Delete("/:id", user.DeleteUser)

	quizEndpoint := endpoint.Group("/quiz", logger.New())
	quizEndpoint.Get("/", quiz.GetQuizzes)
	quizEndpoint.Get("/:id", quiz.GetQuiz)
	quizEndpoint.Post("/", quiz.NewQuiz)
	quizEndpoint.Put("/", quiz.UpdateQuiz)
	quizEndpoint.Delete("/:id", quiz.DeleteQuiz)

	quizResultEndpoint := endpoint.Group("/quizResult", logger.New())
	quizResultEndpoint.Get("/", quiz.GetQuizResults)
	quizResultEndpoint.Get("/:id", quiz.GetQuizResult)
	quizResultEndpoint.Post("/", quiz.NewQuizResult)
	quizResultEndpoint.Put("/", quiz.UpdateQuizResult)
	quizResultEndpoint.Delete("/:id", quiz.DeleteQuizResult)

	teacherEndpoint := endpoint.Group("/teacher", logger.New())
	teacherEndpoint.Get("/", teacher.GetTeachers)
	teacherEndpoint.Get("/:id", teacher.GetTeacher)
	teacherEndpoint.Post("/", teacher.NewTeacher)
	teacherEndpoint.Put("/", teacher.UpdateTeacher)
	teacherEndpoint.Delete("/:id", teacher.DeleteTeacher)

	studentEndpoint := endpoint.Group("/student", logger.New())
	studentEndpoint.Get("/", student.GetStudents)
	studentEndpoint.Get("/:id", student.GetStudent)
	studentEndpoint.Post("/", student.NewStudent)
	studentEndpoint.Put("/", student.UpdateStudent)
	studentEndpoint.Delete("/:id", student.DeleteStudent)

	subjectEndpoint := endpoint.Group("/subject", logger.New())
	subjectEndpoint.Get("/", subject.GetSubjects)
	subjectEndpoint.Get("/:id", subject.GetSubject)
	subjectEndpoint.Post("/", subject.NewSubject)
	subjectEndpoint.Put("/", subject.UpdateSubject)
	subjectEndpoint.Delete("/:id", subject.DeleteSubject)

	sectionEndpoint := endpoint.Group("/section", logger.New())
	sectionEndpoint.Get("/", section.GetSections)
	sectionEndpoint.Get("/:id", section.GetSection)
	sectionEndpoint.Post("/", section.NewSection)
	sectionEndpoint.Put("/", section.UpdateSection)
	sectionEndpoint.Delete("/:id", section.DeleteSection)

	quizResultEndpoint.Post("/check", quiz.CheckQuizResult)

	err = app.Listen(fmt.Sprintf(":%s", envRouting.Port))

	if err != nil {
		log.Fatal(err.Error())
	}
}
