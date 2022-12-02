package ticker

import (
	"github.com/Niexiawei/logrotate"
	"time"
)

func CalRotateTimeDuration(now time.Time, duration time.Duration) time.Duration {
	nowUnixNao := now.UnixNano()
	NanoSecond := duration.Nanoseconds()
	nextRotateTime := NanoSecond - (nowUnixNao % NanoSecond)
	return time.Duration(nextRotateTime)
}

func SetTime(date logrotate.TimeTickerDate) (d time.Duration) {
	now := time.Now()
	setTime := time.Date(now.Year(), now.Month(), now.Day(), date.Hour, date.Min, date.Sec, 0, now.Location())
	d = setTime.Sub(now)
	if d > 0 {
		return
	}
	return d + time.Hour*24
}
