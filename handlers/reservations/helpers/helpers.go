package helpers

import (
	"combustiblemon/keletron-tennis-be/database/models/CourtModel"
	"combustiblemon/keletron-tennis-be/database/models/ReservationModel"
	"combustiblemon/keletron-tennis-be/modules/helpers"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

const TIME_SPLIT string = "T"

func AddMinutesToTime(t string, minutes int) string {
	loc, _ := time.LoadLocation("Europe/Athens")

	hour, _ := strconv.Atoi(strings.Split(t, ":")[0])
	minute, _ := strconv.Atoi(strings.Split(t, ":")[1])

	tm := time.Date(1970, time.Month(1), 1, hour, minute, 0, 0, loc)
	tm = tm.Add(time.Minute * time.Duration(minutes))

	return fmt.Sprintf("%02d:%02d", tm.Hour(), tm.Minute())
}

type OverlappingTimeData struct {
	StartTime string
	EndTime   string
	Duration  int
}

func IsTimeOverlapping(
	reservation OverlappingTimeData,
	against OverlappingTimeData,
) bool {
	// if start time is within the reservation time
	return (against.StartTime < reservation.StartTime &&
		reservation.StartTime < against.EndTime) ||
		// or end time is within the reservation time
		(against.StartTime < reservation.EndTime &&
			reservation.EndTime < against.EndTime) ||
		// or if the times are the same
		(against.StartTime == reservation.StartTime &&
			reservation.Duration == against.Duration)
}

func IsReservationTimeFree(
	courtReservations []ReservationModel.Reservation,
	courtReservedTimes []CourtModel.ReservedTimes,
	datetime string,
	duration int,
	reservationID string,
) bool {
	reservationCheck := true

	startTime := strings.Split(datetime, TIME_SPLIT)[1]
	endTime := AddMinutesToTime(startTime, duration)

	if len(courtReservations) > 0 {
		reservationsToCheck := []ReservationModel.Reservation{}

		for _, r := range courtReservations {
			dateCheck := strings.Split(r.Datetime, TIME_SPLIT)[0] == strings.Split(datetime, TIME_SPLIT)[0]

			if reservationID != "" {
				dateCheck = r.ID.String() == reservationID && dateCheck
			}

			if dateCheck {
				reservationsToCheck = append(reservationsToCheck, r)
			}
		}

		if len(reservationsToCheck) == 0 {
			reservationCheck = true
		} else {
		out:
			for _, r := range reservationsToCheck {
				rstartTime := strings.Split(r.Datetime, TIME_SPLIT)[1]

				if (IsTimeOverlapping(OverlappingTimeData{
					StartTime: startTime,
					EndTime:   endTime,
					Duration:  duration,
				}, OverlappingTimeData{
					Duration:  r.Duration,
					EndTime:   AddMinutesToTime(rstartTime, r.Duration),
					StartTime: rstartTime,
				})) {
					reservationCheck = false
					break out
				}
			}
		}
	}

	if len((courtReservedTimes)) == 0 || !reservationCheck {
		return reservationCheck
	}

	weekDay := strings.ToUpper(helpers.ParseDate(datetime).Weekday().String())

	reservedTimesCheck := true

out2:
	for _, r := range courtReservedTimes {
		if !slices.Contains(r.Days, weekDay) {
			continue
		}

		if (IsTimeOverlapping(
			OverlappingTimeData{
				startTime,
				endTime,
				duration,
			},
			OverlappingTimeData{
				Duration:  r.Duration,
				EndTime:   AddMinutesToTime(r.StartTime, r.Duration),
				StartTime: r.StartTime,
			},
		)) {
			reservedTimesCheck = false
			break out2
		}
	}

	return reservationCheck && reservedTimesCheck
}
