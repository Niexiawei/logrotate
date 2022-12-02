package ticker

import (
	"testing"
	"time"
)

func TestSetTime(t *testing.T) {
	timeD := SetTime(TimeTickerDate{})
	t.Log(time.Now().Add(timeD).Format("2006-01-02 15:04:05"))
}
