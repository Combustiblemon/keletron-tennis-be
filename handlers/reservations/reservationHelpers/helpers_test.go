package reservationHelpers

import (
	"testing"
)

func TestAddMinutesToTime(t *testing.T) {
	time := AddMinutesToTime("09:32", 10)

	if time != "09:42" {
		t.Fatalf(`error`)
	}

	time = AddMinutesToTime("09:55", 10)

	if time != "10:05" {
		t.Fatalf(`error`)
	}

	time = AddMinutesToTime("23:59", 10)

	if time != "00:09" {
		t.Fatalf(`error`)
	}
}

func TestIsTimeOverLapping(t *testing.T) {
	res := IsTimeOverlapping(OverlappingTimeData{
		StartTime: "10:00",
		EndTime:   "11:30",
		Duration:  90,
	},
		OverlappingTimeData{
			StartTime: "12:00",
			EndTime:   "13:30",
			Duration:  90,
		})

	if res == true {
		t.Fatalf(`error`)
	}
}
