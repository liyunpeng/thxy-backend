package utils

import (
	"strconv"
	"time"
)

func GetTimeZeroLastSecond(timeDuration string) string {
	now := time.Now()
	duration, _ := time.ParseDuration(timeDuration)
	newTime := now.Add(duration)
	newTimeStr := newTime.Format("2006-01-02") + " 00:00:00"
	return newTimeStr
}

func GetTimeStrFromSecond(seconds int) string {
	minute := strconv.Itoa(seconds / 60)
	second := strconv.Itoa(seconds % 60)

	if second == "0" || second == "" {
		return minute + ":00"
	} else {
		return minute + ":" + second
	}
}
