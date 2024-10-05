package announcements

import (
	AnnouncementModel "combustiblemon/keletron-tennis-be/database/models/announcement"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		announcements, err := AnnouncementModel.Find(bson.D{{Key: "visible", Value: true}})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, map[string]any{})

			return
		}

		ctx.JSON(http.StatusOK, announcements)
	}
}
