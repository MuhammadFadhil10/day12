package helper

import (
	"strconv"
	"time"
)

func GetDuration(startDate time.Time, endDate time.Time) string {
	var duration string

	var startMs = startDate.UnixMicro()
	var endMs = endDate.UnixMicro()

	margin := ((endMs - startMs) / (1000 * 60 * 60 * 24) / 1000)

	// fmt.Println(margin)

	if margin < 30 {
		if margin == 0 {
			duration = "a few hours";
		} else {
			duration = strconv.Itoa(int(margin)) + " Day"
		}
	}  else {
		if margin < 365 {
			if margin % 30 == 0 {
				duration = strconv.Itoa(int(margin / 30)) + " Month"
			} else {
				duration = strconv.Itoa(int(margin / 30)) + " Month " + strconv.Itoa(int(margin % 30)) + " Day"
			}
		} else {
			duration = strconv.Itoa(int(margin / 365)) + " Year"
		}
	}

	return duration
}