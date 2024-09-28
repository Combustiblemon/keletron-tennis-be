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
	"combustiblemon/keletron-tennis-be/handlers/users"
	"combustiblemon/keletron-tennis-be/middleware"
	"log"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const SERVER_PORT = "2000"

// // album represents data about a record album.
// type album struct {
// 	ID     string  `json:"id"`
// 	Title  string  `json:"title"`
// 	Artist string  `json:"artist"`
// 	Price  float64 `json:"price"`
// }

// // albums slice to seed record album data.
// var albums = []album{
// 	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
// 	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
// 	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
// }

// // getAlbums responds with the list of all albums as JSON.
// func getAlbums(c *gin.Context) {
// 	c.JSON(http.StatusOK, albums)
// }

// // postAlbums adds an album from JSON received in the request body.
// func postAlbums(c *gin.Context) {
// 	var newAlbum album

// 	// Call BindJSON to bind the received JSON to
// 	// newAlbum.
// 	if err := c.BindJSON(&newAlbum); err != nil {
// 		return
// 	}

// 	// Add the new album to the slice.
// 	albums = append(albums, newAlbum)
// 	c.JSON(http.StatusCreated, newAlbum)
// }

// // getAlbumByID locates the album whose ID value matches the id
// // parameter sent by the client, then returns that album as a response.
// func getAlbumByID(c *gin.Context) {
// 	id := c.Param("id")

// 	// Loop over the list of albums, looking for
// 	// an album whose ID value matches the parameter.
// 	for _, a := range albums {
// 		if a.ID == id {
// 			c.JSON(http.StatusOK, a)
// 			return
// 		}
// 	}

// 	c.JSON(http.StatusNotFound, gin.H{"message": "album not found"})
// }

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
			reservationsGroup.GET("/", reservations.GET())
			reservationsGroup.POST("/", reservations.POST())
			reservationsGroup.GET("/:id", reservations.GET_ID())
			reservationsGroup.PUT("/:id", reservations.PUT_ID())
			reservationsGroup.DELETE("/:id", reservations.DELETE_ID())
		}

		courtsGroup := authorized.Group("courts")
		{
			courtsGroup.GET("/", courts.GET())
			courtsGroup.GET("/:id", courts.GET_ID())
		}

		usersGroup := authorized.Group("users")
		{
			usersGroup.GET("/", users.GET())
			usersGroup.PUT("/", users.PUT())
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

//revive:enable:add-constant

func main() {
	if err := godotenv.Load("secret.env", ".env"); err != nil {
		log.Println("No .env file found")
	}

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

	router.Use(middleware.Logger())

	setupAuthGroup(router)

	setupAuthorizedGroup(router)

	router.Use(middleware.Error())

	err = router.Run(fmt.Sprintf("localhost:%v", SERVER_PORT))

	if err != nil {
		log.Fatal("Error bringing server online", err)
	}
}
