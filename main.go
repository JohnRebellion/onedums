package main

import (
	"fmt"
	"log"
	"onedums/anouncement"
	"onedums/envRouting"
	"onedums/learningMaterial"
	"onedums/quiz"
	"onedums/section"
	"onedums/student"
	"onedums/subject"
	"onedums/teacher"
	"onedums/twilioService"
	"onedums/user"
	"os"
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
	makeDirectoryIfNotExists("files/learningMaterials")
	twilioService.NewClient(envRouting.TwilioAccountSID, envRouting.TwilioAuthenticationToken, envRouting.TwilioPhoneNumber)
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
		&anouncement.Anouncement{},
		&learningMaterial.LearningMaterial{},
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
	app.Use(logger.New())
	app.Static("/", envRouting.StaticWebLocation)
	app.Static("/learningMaterials", "files/learningMaterials")
	setupPublicRoutes(app)
	setupPrivateRoutes(app)

	err = app.Listen(fmt.Sprintf(":%s", envRouting.Port))

	if err != nil {
		log.Fatal(err.Error())
	}
}

func setupPublicRoutes(app *fiber.App) {
	apiEndpoint := app.Group("/api")
	v1Endpoint := apiEndpoint.Group("/v1")
	userInfoEndpoint := v1Endpoint.Group("/userInfo")
	userInfoEndpoint.Post("/", user.NewUserInfo)
	userInfoEndpoint.Post("/auth", user.Authenticate)
}

func setupPrivateRoutes(app *fiber.App) {
	app.Use(fiberUtils.AuthenticationMiddleware(fiberUtils.JWTConfig{
		Duration:     15 * time.Minute,
		CookieMaxAge: 15 * 60,
		SetCookies:   true,
		SecretKey:    []byte(envRouting.SecretKey),
	}))
	apiEndpoint := app.Group("/api")
	v1Endpoint := apiEndpoint.Group("/v1")
	userInfoEndpoint := v1Endpoint.Group("/userInfo")

	userInfoEndpoint.Get("/:id", user.GetUserInfo)
	userInfoEndpoint.Get("/", user.GetUserInfos)
	userInfoEndpoint.Put("/", user.UpdateUserInfo)
	userInfoEndpoint.Delete("/:id", user.DeleteUserInfo)

	userEndpoint := v1Endpoint.Group("/user")
	userEndpoint.Get("/", user.GetUsers)
	userEndpoint.Get("/:id", user.GetUser)
	userEndpoint.Post("/", user.NewUser)
	userEndpoint.Put("/", user.UpdateUser)
	userEndpoint.Delete("/:id", user.DeleteUser)

	quizEndpoint := v1Endpoint.Group("/quiz")
	quizEndpoint.Get("/", quiz.GetQuizzes)
	quizEndpoint.Get("/:id", quiz.GetQuiz)
	quizEndpoint.Post("/", quiz.NewQuiz)
	quizEndpoint.Put("/", quiz.UpdateQuiz)
	quizEndpoint.Delete("/:id", quiz.DeleteQuiz)

	quizResultEndpoint := v1Endpoint.Group("/quizResult")
	quizResultEndpoint.Get("/", quiz.GetQuizResults)
	quizResultEndpoint.Get("/:id", quiz.GetQuizResult)
	quizResultEndpoint.Post("/", quiz.NewQuizResult)
	quizResultEndpoint.Put("/", quiz.UpdateQuizResult)
	quizResultEndpoint.Delete("/:id", quiz.DeleteQuizResult)

	teacherEndpoint := v1Endpoint.Group("/teacher")
	teacherEndpoint.Get("/", teacher.GetTeachers)
	teacherEndpoint.Get("/:id", teacher.GetTeacher)
	teacherEndpoint.Post("/", teacher.NewTeacher)
	teacherEndpoint.Put("/", teacher.UpdateTeacher)
	teacherEndpoint.Delete("/:id", teacher.DeleteTeacher)

	studentEndpoint := v1Endpoint.Group("/student")
	studentEndpoint.Get("/", student.GetStudents)
	studentEndpoint.Get("/:id", student.GetStudent)
	studentEndpoint.Post("/", student.NewStudent)
	studentEndpoint.Put("/", student.UpdateStudent)
	studentEndpoint.Delete("/:id", student.DeleteStudent)

	subjectEndpoint := v1Endpoint.Group("/subject")
	subjectEndpoint.Get("/", subject.GetSubjects)
	subjectEndpoint.Get("/:id", subject.GetSubject)
	subjectEndpoint.Post("/", subject.NewSubject)
	subjectEndpoint.Put("/", subject.UpdateSubject)
	subjectEndpoint.Delete("/:id", subject.DeleteSubject)

	sectionEndpoint := v1Endpoint.Group("/section")
	sectionEndpoint.Get("/", section.GetSections)
	sectionEndpoint.Get("/:id", section.GetSection)
	sectionEndpoint.Post("/", section.NewSection)
	sectionEndpoint.Put("/", section.UpdateSection)
	sectionEndpoint.Delete("/:id", section.DeleteSection)

	quizResultEndpoint.Post("/check", quiz.CheckQuizResult)

	anouncementEndpoint := v1Endpoint.Group("/anouncement")
	anouncementEndpoint.Get("/", anouncement.GetAnouncements)
	anouncementEndpoint.Get("/:id", anouncement.GetAnouncement)
	anouncementEndpoint.Post("/", anouncement.NewAnouncement)
	anouncementEndpoint.Put("/", anouncement.UpdateAnouncement)
	anouncementEndpoint.Delete("/:id", anouncement.DeleteAnouncement)

	learningMaterialEndpoint := v1Endpoint.Group("/learningMaterial")
	learningMaterialCurrentUserEndpoint := learningMaterialEndpoint.Group("/currentUser")
	learningMaterialCurrentUserEndpoint.Get("/", learningMaterial.GetLearningMaterialsCurrentUser)
	learningMaterialCurrentUserEndpoint.Get("/:id", learningMaterial.GetLearningMaterialCurrentUser)
	learningMaterialCurrentUserEndpoint.Put("/", learningMaterial.UpdateLearningMaterialCurrentUser)
	learningMaterialCurrentUserEndpoint.Delete("/:id", learningMaterial.DeleteLearningMaterialCurrentUser)
	learningMaterialCurrentUserEndpoint.Post("/uploadFile", learningMaterial.UploadLearningMaterialCurrentUser)
	learningMaterialCurrentUserEndpoint.Put("/uploadFile", learningMaterial.UploadUpdatedLearningMaterialCurrentUser)

	learningMaterialEndpoint.Get("/", learningMaterial.GetLearningMaterials)
	learningMaterialEndpoint.Get("/:id", learningMaterial.GetLearningMaterial)
	learningMaterialEndpoint.Post("/", learningMaterial.NewLearningMaterial)
	learningMaterialEndpoint.Put("/", learningMaterial.UpdateLearningMaterial)
	learningMaterialEndpoint.Delete("/:id", learningMaterial.DeleteLearningMaterial)
	learningMaterialEndpoint.Post("/uploadFile", learningMaterial.UploadLearningMaterial)
	learningMaterialEndpoint.Get("/:id/downloadFile", learningMaterial.DownloadLearningMaterialFile)
	learningMaterialEndpoint.Put("/uploadFile", learningMaterial.UploadUpdatedLearningMaterial)

	app.Static("/learningMaterials", "files/learningMaterials")
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModeDir|0755)
	}

	return nil
}
