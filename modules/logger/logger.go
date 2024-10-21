package logger

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Info(ctx *gin.Context, msg string) {
	info, _ := ctx.Get("info")

	// Print the log message in the same format as Gin
	slog.Info("[INFO] %v              | %15s |%-7s %s | Info: %s | %s\n",
		time.Now().Format(time.RFC1123),
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		info,
		msg,
	)
}

func Error(ctx *gin.Context, msg string) {
	info, _ := ctx.Get("info")

	// Print the log message in the same format as Gin
	slog.Info("[ERROR] %v              | %15s |%-7s %s | Info: %s | %s\n",
		time.Now().Format(time.RFC1123),
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		info,
		msg,
	)
}

func Warn(ctx *gin.Context, msg string) {
	info, _ := ctx.Get("info")

	// Print the log message in the same format as Gin
	slog.Info("[Warn] %v              | %15s |%-7s %s | Info: %s | %s\n",
		time.Now().Format(time.RFC1123),
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		info,
		msg,
	)
}
