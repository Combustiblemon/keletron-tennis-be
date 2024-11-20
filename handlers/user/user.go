package user

import (
	"combustiblemon/keletron-tennis-be/modules/errorHandler"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"combustiblemon/keletron-tennis-be/modules/logger"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user, exists := helpers.GetUser(ctx)

		if exists {
			ctx.JSON(http.StatusOK, user.Sanitize())
			return
		}

		errorHandler.SendError(ctx, http.StatusNotFound, fmt.Errorf("userNotFound"))
	}
}

func PutOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user, exists := helpers.GetUser(ctx)

		if exists {
			bodyAsByteArray, err := io.ReadAll(ctx.Request.Body)

			if err != nil {
				ctx.Status(http.StatusBadRequest)
				logger.Debug(ctx, err.Error())
				return
			}

			var data map[string]any
			err = json.Unmarshal(bodyAsByteArray, &data)

			if err != nil {
				errorHandler.SendError(ctx, http.StatusBadRequest, err)
				return
			}

			if data["FCMToken"] == nil && data["name"] == nil {
				errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("no data received"))
				return
			}

			name := data["FCMToken"]

			if name != nil {
				if nameStr, ok := name.(string); ok {
					user.Name = nameStr
				} else {
					errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("invalid type for name"))
					return
				}
			}

			user.Save()

			return
		}

		errorHandler.SendError(ctx, http.StatusNotFound, fmt.Errorf("userNotFound"))
	}
}
