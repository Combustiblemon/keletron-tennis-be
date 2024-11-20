package reservations

import (
	"combustiblemon/keletron-tennis-be/database/models/CourtModel"
	"combustiblemon/keletron-tennis-be/database/models/ReservationModel"
	resHelpers "combustiblemon/keletron-tennis-be/handlers/reservations/reservationHelpers"
	"combustiblemon/keletron-tennis-be/modules/errorHandler"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"combustiblemon/keletron-tennis-be/modules/logger"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/mold/v4/scrubbers"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
)

func GetOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_id := ctx.Query(("id"))

		if _id == "" {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("no id provided"))
			return
		}

		reservation, err := ReservationModel.FindOne(bson.D{{Key: "_id", Value: _id}})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		user, exists := helpers.GetUser(ctx)

		if exists && reservation.Owner.String() == user.ID.String() {
			ctx.JSON(http.StatusOK, reservation.SanitizeOwner())
		}

		ctx.JSON(http.StatusOK, reservation.Sanitize())
	}
}

func DeleteOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_id := ctx.Query(("id"))

		if _id == "" {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("no id provided"))
			return
		}

		reservation, err := ReservationModel.FindOne(bson.D{{Key: "_id", Value: _id}})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		user, exists := helpers.GetUser(ctx)

		if exists && reservation.Owner.String() != user.ID.String() {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		err = reservation.Delete()

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
		}
	}
}

var (
	conform  = modifiers.New()
	scrub    = scrubbers.New()
	validate = validator.New()
)

func PutOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyAsByteArray, err := io.ReadAll(ctx.Request.Body)

		if err != nil {
			ctx.Status(http.StatusBadRequest)
			logger.Debug(ctx, err.Error())
			return
		}

		var r ReservationModel.Reservation
		err = json.Unmarshal(bodyAsByteArray, &r)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusBadRequest, err)
			return
		}

		err = conform.Struct(context.Background(), &r)
		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		err = validate.Struct(r)
		if err != nil {
			if validatorErr := (validator.ValidationErrors{}); errors.As(err, &validatorErr) {
				errorHandler.SendError(ctx, http.StatusBadRequest, err)
			} else {
				errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			}

			return
		}

		if !resHelpers.IsTimeValid(r.Datetime) {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("invalidDatetime"))

			return
		}

		rOld, err := ReservationModel.FindOne(bson.D{{Key: "_id", Value: r.ID}})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		user, exists := helpers.GetUser(ctx)

		if !exists {
			slog.Error("No ctx user found in reservation.PostOne")
			ctx.Status(http.StatusInternalServerError)
			return
		}

		if user.ID.String() != rOld.Owner.String() {
			ctx.Status(http.StatusUnauthorized)
			return
		}

		court, err := CourtModel.FindOne(bson.D{{Key: "_id", Value: rOld.Court}})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		if court == nil {
			errorHandler.SendError(ctx, http.StatusNotFound, fmt.Errorf("court.notFound"))

			return
		}

		reservations, err := ReservationModel.Find(bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$ne", Value: rOld.ID},
			}},
			{Key: "court", Value: rOld.Court},
			{Key: "datetime",
				Value: bson.D{
					{Key: "$gte", Value: "2024-08-04T09:00"},
					{Key: "$lte", Value: "2024-08-05T13:00"},
				},
			},
		})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)

			return
		}

		if !resHelpers.IsReservationTimeFree(*reservations, court.ReservationsInfo.ReservedTimes, r.Datetime, court.ReservationsInfo.Duration, "") {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("reservationExistsOnTime"))
			return
		}

		err = rOld.Save(&r)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, r)
	}
}

func PostOne() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bodyAsByteArray, err := io.ReadAll(ctx.Request.Body)

		if err != nil {
			ctx.Status(http.StatusBadRequest)
			logger.Debug(ctx, err.Error())
			return
		}

		var r ReservationModel.Reservation
		err = json.Unmarshal(bodyAsByteArray, &r)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusBadRequest, err)
			return
		}

		err = conform.Struct(context.Background(), &r)
		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		err = validate.Struct(r)
		if err != nil {
			if validatorErr := (validator.ValidationErrors{}); errors.As(err, &validatorErr) {
				errorHandler.SendError(ctx, http.StatusBadRequest, err)
			} else {
				errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			}

			return
		}

		if !resHelpers.IsTimeValid(r.Datetime) {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("invalidDatetime"))

			return
		}

		court, err := CourtModel.FindOne(bson.D{{Key: "_id", Value: r.Court}})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)

			return
		}

		if court == nil {
			errorHandler.SendError(ctx, http.StatusNotFound, fmt.Errorf("court.notFound"))

			return
		}

		reservations, err := ReservationModel.Find(bson.D{
			{Key: "court", Value: r.Court},
			{Key: "datetime",
				Value: bson.D{
					{Key: "$gte", Value: "2024-08-04T09:00"},
					{Key: "$lte", Value: "2024-08-05T13:00"},
				},
			},
		})

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)

			return
		}

		user, exists := helpers.GetUser(ctx)

		if !exists {
			slog.Error("No ctx user found in reservation.PostOne")
			ctx.Status(http.StatusInternalServerError)
			return
		}

		r.Owner = user.ID

		if !resHelpers.IsReservationTimeFree(*reservations, court.ReservationsInfo.ReservedTimes, r.Datetime, court.ReservationsInfo.Duration, "") {
			errorHandler.SendError(ctx, http.StatusBadRequest, fmt.Errorf("reservationExistsOnTime"))
			return
		}

		rNew, err := ReservationModel.Create(&r)

		if err != nil {
			errorHandler.SendError(ctx, http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusCreated, rNew)
	}
}
