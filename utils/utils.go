package utils

import "time"

func GetTimeZeroLastSecond( timeDuration string) string{
	now := time.Now()
	duration, _ := time.ParseDuration(timeDuration)
	newTime := now.Add(duration)
	newTimeStr := newTime.Format("2006-01-02") + " 00:00:00"
	return newTimeStr
}