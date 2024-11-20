package main

import (
	"combustiblemon/keletron-tennis-be/database"
	"combustiblemon/keletron-tennis-be/handlers/admin/adminCourts"
	"combustiblemon/keletron-tennis-be/handlers/admin/adminReservations"
	"combustiblemon/keletron-tennis-be/handlers/admin/adminUsers"
	"combustiblemon/keletron-tennis-be/handlers/auth"
	"combustiblemon/keletron-tennis-be/handlers/auth/providersGoogle"
	"combustiblemon/keletron-tennis-be/handlers/courts"
	"combustiblemon/keletron-tennis-be/handlers/reservations"
	"combustiblemon/keletron-tennis-be/handlers/user"
	"combustiblemon/keletron-tennis-be/middleware"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"log"
	"os"
	"time"

	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var serverPort = helpers.Condition(os.Getenv("PORT") == "", "2000", os.Getenv("PORT"))

func setupAuthGroup(router *gin.Engine) {
	authGroup := router.Group("auth")
	{
		authGroup.GET("/session", auth.Session())
		authGroup.POST("/login", auth.Login())
		authGroup.POST("/register", auth.Register())

		callbackGroup := authGroup.Group("providers")

		{
			callbackGroup.GET("/google/start", providersGoogle.Start())
			callbackGroup.GET("/google/callback", providersGoogle.Callback())
		}
	}
}

//revive:disable:add-constant
func setupAuthorizedGroup(router *gin.Engine) {
	authorized := router.Group("/")

	authorized.Use(middleware.Auth())
	{
		reservationsGroup := authorized.Group("reservations")
		{
			reservationsGroup.GET("/", reservations.GetMany())
			reservationsGroup.POST("/", reservations.PostOne())
			reservationsGroup.GET("/:id", reservations.GetOne())
			reservationsGroup.PUT("/:id", reservations.PutOne())
			reservationsGroup.DELETE("/:id", reservations.DeleteOne())
		}

		courtsGroup := authorized.Group("courts")
		{
			courtsGroup.GET("/", courts.GetMany())
			courtsGroup.GET("/:id", courts.GetOne())
		}

		usersGroup := authorized.Group("user")
		{
			usersGroup.GET("/", user.GetOne())
			usersGroup.PUT("/", user.PutOne())
		}

		admin := authorized.Group("admin")

		admin.Use(middleware.Admin())
		{
			reservationsGroup := admin.Group("reservations")
			{
				reservationsGroup.GET("/", adminReservations.GET())
				reservationsGroup.POST("/", adminReservations.POST())
				reservationsGroup.GET("/:id", adminReservations.GET_ID())
				reservationsGroup.PUT("/:id", adminReservations.PUT_ID())
				reservationsGroup.DELETE("/:id", adminReservations.DELETE_ID())
			}

			courtsGroup := admin.Group("courts")
			{
				courtsGroup.GET("/", adminCourts.GET())
				courtsGroup.POST("/", adminCourts.POST())
				courtsGroup.GET("/:id", adminCourts.GET_ID())
				courtsGroup.PUT("/:id", adminCourts.PUT_ID())
				courtsGroup.DELETE("/:id", adminCourts.DELETE_ID())
			}

			usersGroup := admin.Group("users")
			{
				usersGroup.GET("/", adminUsers.GET())
				usersGroup.PUT("/", adminUsers.PUT())
				usersGroup.GET("/:id", adminUsers.GET_ID())
				usersGroup.POST("/:id", adminUsers.POST_ID())
				usersGroup.PUT("/:id", adminUsers.PUT_ID())
				usersGroup.DELETE("/:id", adminUsers.DELETE_ID())
			}
		}
	}
}

func main() {
	if err := godotenv.Load("secret.env", ".env"); err != nil {
		log.Println("No .env file found")
	}

	providersGoogle.Init()

	err := database.Setup()

	// nolint:errcheck
	defer database.Teardown()

	if err != nil {
		log.Fatalln("Error setting up database\n", err)
	}

	router := gin.Default()
	err = router.SetTrustedProxies(nil)

	if err != nil {
		log.Fatal("error in SetTrustedProxies", err)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:2000"},
		AllowMethods:     []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "Content-Length", "X-CSRF-Token", "Token", "session", "Origin", "Host", "Connection", "Accept-Encoding", "Accept-Language", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		MaxAge: 24 * time.Hour,
	}))

	router.Use(middleware.Info())
	router.Use(middleware.Logger())

	router.GET("/announcements")

	setupAuthGroup(router)

	setupAuthorizedGroup(router)

	router.Use(middleware.Error())

	err = router.Run(fmt.Sprintf("localhost:%v", serverPort))

	if err != nil {
		log.Fatal("Error bringing server online", err)
	}
}
