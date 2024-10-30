package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Info() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyAsByteArray, err := io.ReadAll(ctx.Request.Body)

		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}

		jsonMap := map[string]any{}
		err = json.Unmarshal(bodyAsByteArray, &jsonMap)

		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
		}

		ctx.Set("info", time.Now().UnixNano())
		ctx.Set("json", jsonMap)

		ctx.Next()
	}
}
