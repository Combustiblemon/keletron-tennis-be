package courts

import (
	"combustiblemon/keletron-tennis-be/database/models/CourtModel"
	"combustiblemon/keletron-tennis-be/modules/errorHandler"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMany() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		court, err := CourtModel.FindOne(bson.D{})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, court)
	}
}
