package ticker

import (
	"testing"
	"time"
)

func Test_Output(t *testing.T) {
	d := CalRotateTimeDuration(time.Now(), time.Hour*2)
	t.Logf("now:%s, duration:%s\n", time.Now(), d)
}
