package logger

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func Info(ctx *gin.Context, msg string) {
	info, _ := ctx.Get("info")

	// Print the log message in the same format as Gin
	slog.Info(fmt.Sprintf("[INFO] %s |%s %s | Info: %s | %s\n",
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		info,
		msg,
	))
}

func Error(ctx *gin.Context, msg string) {
	info, _ := ctx.Get("info")

	// Print the log message in the same format as Gin
	slog.Error(fmt.Sprintf("[ERROR] %s | %s %s | Info: %s | %s\n",
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		info,
		msg,
	))
}

func Warn(ctx *gin.Context, msg string) {
	info, _ := ctx.Get("info")

	// Print the log message in the same format as Gin
	slog.Warn(fmt.Sprintf("[Warn] %s | %s %s | Info: %s | %s\n",
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		info,
		msg,
	))
}

func Debug(ctx *gin.Context, msg string) {
	if gin.Mode() == gin.ReleaseMode {
		return
	}

	slog.Info(fmt.Sprintf("[Debug] %s | %s %s | %s\n",
		ctx.ClientIP(),
		ctx.Request.Method,
		ctx.Request.URL.Path,
		msg,
	))
}
