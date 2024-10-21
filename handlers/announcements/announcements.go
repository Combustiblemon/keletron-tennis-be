package announcements

import (
	"combustiblemon/keletron-tennis-be/database/models/AnnouncementModel"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		announcements, err := AnnouncementModel.Find(bson.D{{Key: "visible", Value: true}})

		if err != nil {
			helpers.SendError(ctx, http.StatusInternalServerError, err)

			return
		}

		ctx.JSON(http.StatusOK, announcements)
	}
}
