package errorHandler

import (
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"combustiblemon/keletron-tennis-be/modules/logger"
	"log/slog"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ValidationErrorInfo struct {
	Key   string
	Error string
	Info  string
}

func GenerateValidationError(err error) []ValidationErrorInfo {
	errParts := strings.Split(err.Error(), "\n")
	regKey, _ := regexp2.Compile("(?<=Key: ')\\b\\w+\\.\\w+\\b", 0)
	regTag, _ := regexp2.Compile("(?<=failed on the ')\\b\\w+\\b(?=' tag)", 0)

	errs := []ValidationErrorInfo{}

	for _, v := range errParts {
		if strings.Contains(v, "Error:Field validation") {
			m1, _ := regKey.FindStringMatch(v)
			m2, _ := regTag.FindStringMatch(v)

			key := helpers.FirstToLower(strings.Replace(m1.String(), "ID", "_id", 1))

			tag := strings.ToLower(m2.String())

			errs = append(errs, ValidationErrorInfo{
				Key:   key,
				Error: tag,
				Info:  "",
			})
		}
	}

	return errs
}

func ObjectIDFromHex(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)

	if err != nil {
		slog.Error(err.Error())
		return primitive.NilObjectID
	}

	return id
}

func FormatUmarshalError(err error) ValidationErrorInfo {
	if strings.Contains(strings.ToLower(err.Error()), "objectid") {
		return ValidationErrorInfo{}
	}

	return ValidationErrorInfo{}
}

func SendError(ctx *gin.Context, status int, err error) {
	logger.Warn(ctx, err.Error())

	if strings.Contains(err.Error(), "Error:Field validation") {
		ctx.JSON(status, map[string]any{
			"errors": GenerateValidationError(err),
		})

		return
	}

	ctx.JSON(status, map[string]any{
		"errors": []string{err.Error()},
	})
}
