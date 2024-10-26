package courts

import (
	"combustiblemon/keletron-tennis-be/database/models/CourtModel"
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

		court, err := CourtModel.FindOne(bson.D{{Key: "_id", Value: _id}})

		if err != nil {
			helpers.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, court)
	}
}
