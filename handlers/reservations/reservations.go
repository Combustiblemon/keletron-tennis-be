package reservations

import (
	"combustiblemon/keletron-tennis-be/database/models/ReservationModel"
	"combustiblemon/keletron-tennis-be/modules/errorHandler"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetMany() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, exists := helpers.GetUser(ctx)

		if !exists {
			errorHandler.SendError(ctx, http.StatusInternalServerError, fmt.Errorf("no user found"))
			return
		}

		reservations, err := ReservationModel.Find(bson.D{{Key: "owner", Value: user.ID}})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		var ret []ReservationModel.ReservationSanitizedOwner

		for _, r := range *reservations {
			ret = append(ret, r.SanitizeOwner())
		}

		ctx.JSON(http.StatusOK, ret)
	}
}

func DeleteMany() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
