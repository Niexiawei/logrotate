package ticker

import (
	"time"
)

type TimeTickerDate struct {
	Hour int
	Min  int
	Sec  int
}

func SetTime(date TimeTickerDate) (d time.Duration) {
	now := time.Now()
	setTime := time.Date(now.Year(), now.Month(), now.Day(), date.Hour, date.Min, date.Sec, 0, now.Location())
	d = setTime.Sub(now)
	if d > 0 {
		return
	}
	return d + time.Hour*24
}
