package main

import (
	admin_courts "combustiblemon/keletron-tennis-be/handlers/admin/courts"
	admin_reservations "combustiblemon/keletron-tennis-be/handlers/admin/reservations"
	admin_users "combustiblemon/keletron-tennis-be/handlers/admin/users"
	"combustiblemon/keletron-tennis-be/handlers/courts"
	"combustiblemon/keletron-tennis-be/handlers/reservations"
	"combustiblemon/keletron-tennis-be/handlers/users"
	"combustiblemon/keletron-tennis-be/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

const SERVER_PORT = "8080"

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

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(middleware.Logger())
	// router.Use(middleware.Auth())

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
				reservationsGroup.GET("/", admin_reservations.GET())
				reservationsGroup.POST("/", admin_reservations.POST())
				reservationsGroup.GET("/:id", admin_reservations.GET_ID())
				reservationsGroup.PUT("/:id", admin_reservations.PUT_ID())
				reservationsGroup.DELETE("/:id", admin_reservations.DELETE_ID())
			}

			courtsGroup := admin.Group("courts")
			{
				courtsGroup.GET("/", admin_courts.GET())
				courtsGroup.POST("/", admin_courts.POST())
				courtsGroup.GET("/:id", admin_courts.GET_ID())
				courtsGroup.PUT("/:id", admin_courts.PUT_ID())
				courtsGroup.DELETE("/:id", admin_courts.DELETE_ID())
			}

			usersGroup := admin.Group("users")
			{
				usersGroup.GET("/", admin_users.GET())
				usersGroup.PUT("/", admin_users.PUT())
				usersGroup.GET("/:id", admin_users.GET_ID())
				usersGroup.POST("/:id", admin_users.POST_ID())
				usersGroup.PUT("/:id", admin_users.PUT_ID())
				usersGroup.DELETE("/:id", admin_users.DELETE_ID())
			}
		}
	}

	router.Use(middleware.Error())

	router.Run(fmt.Sprintf("localhost:%v", SERVER_PORT))
}
