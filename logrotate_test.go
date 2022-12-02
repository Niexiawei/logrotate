package logrotate

import (
	"fmt"
	"github.com/Niexiawei/logrotate/internal/ticker"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_SimpleWrite(t *testing.T) {
	// no rotate
	t.Run("no rotate", func(t *testing.T) {
		rl, err := NewRotateLog("./testdata/test.log.2006010215", WithCurLogLinkname("./testdata/test.log"))
		if assert.Empty(t, err) {
			defer rl.Close()
			rl.Write([]byte("hello, world!"))
			compareFileContent(t, rl.getLatestLogPath(time.Now()), "hello, world!")
			compareFileContent(t, "./testdata/test.log", "hello, world!")
			os.RemoveAll("./testdata/")
		}
	})
	// no rotate, on link
	t.Run("no rotate and link", func(t *testing.T) {
		rl, err := NewRotateLog("./testdata/test.log")
		if assert.Empty(t, err) {
			defer rl.Close()
			rl.Write([]byte("hello, world!"))
			content, err := ioutil.ReadFile(rl.getLatestLogPath(time.Now()))
			if assert.Empty(t, err) {
				assert.Equal(t, string(content), "hello, world!")
				t.Log(string(content))
			}
			os.RemoveAll("./testdata/")
		}
	})
}

func Test_Rotate(t *testing.T) {
	rl, err := NewRotateLog("./testdata/test.log.2006010215", WithRotateTime(time.Hour), WithCurLogLinkname("./testdata/test.log"))
	if assert.Empty(t, err) {
		defer rl.Close()
		rotate := make(chan time.Time, 1)
		rl.rotate = rotate

		rl.Write([]byte("hello, world\n"))

		nextHour := time.Now().Add(time.Hour)
		rotate <- nextHour

		time.Sleep(time.Millisecond * 100)
		rl.Write([]byte("hello, world2\n"))
		compareFileContent(t, rl.getLatestLogPath(time.Now()), "hello, world\n")
		compareFileContent(t, rl.getLatestLogPath(nextHour), "hello, world2\n")
		compareFileContent(t, "./testdata/test.log", "hello, world2\n")
		os.RemoveAll("./testdata/")
	}
}

func Test_DeleteExpiredFile(t *testing.T) {
	rl, _ := NewRotateLog("./testdata/test.log.2006010215", WithRotateTime(time.Hour), WithCurLogLinkname("./testdata/test.log"),
		WithDeleteExpiredFile(time.Second, "test.log*"))

	defer func() {
		_ = rl.Close()
	}()
	for i := 0; i < 10; i++ {
		file, _ := os.OpenFile(fmt.Sprintf("./testdata/test.log.%d", i), os.O_CREATE, 0644)
		_ = file.Close()
	}
	matches, _ := filepath.Glob("./testdata/test.log*")
	t.Log(matches)
	time.Sleep(time.Millisecond * 1200)
	rl.rotateFile(time.Now())
	time.Sleep(time.Millisecond * 10)
	matches, _ = filepath.Glob("./testdata/test.log*")
	t.Log(matches)

}

func Test_Speed(t *testing.T) {
	t.Skip()
	rl, err := NewRotateLog("./testdata/test.log", WithRotateTime(time.Hour))
	if assert.Empty(t, err) {
		bg := time.Now()
		for i := 0; i < 1000000; i++ {
			rl.Write([]byte("hello, world\n"))
		}
		t.Log(time.Since(bg))
		os.Remove("./testdata/test.log")
	}
}

func compareFileContent(t *testing.T, filename string, str string) {
	content, err := ioutil.ReadFile(filename)
	t.Log(string(content))
	if assert.Empty(t, err) {
		assert.Equal(t, str, string(content))
	}
}

func Test_afterTime(t *testing.T) {
	now := time.Now()
	t.Log(now.Format("2006-01-02 15:04:05"))
	timeD := ticker.CalRotateTimeDuration(now, 24*time.Hour)
	t.Log(timeD.Hours())
}
