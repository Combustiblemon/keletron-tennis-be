package reservations

import (
	"combustiblemon/keletron-tennis-be/database/models/ReservationModel"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_id := ctx.Query(("id"))

		if _id == "" {
			helpers.SendError(ctx, http.StatusBadRequest, fmt.Errorf("no id provided"))
			return
		}

		reservation, err := ReservationModel.FindOne(bson.D{{Key: "_id", Value: _id}})

		if err != nil {
			helpers.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		user := helpers.GetUser(ctx)

		if user != nil && reservation.Owner.String() == user.ID.String() {

			ctx.JSON(http.StatusOK, reservation.SanitizeOwner())
		}

		ctx.JSON(http.StatusOK, reservation.Sanitize())
	}
}

func PutOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func DeleteOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func PostOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}
