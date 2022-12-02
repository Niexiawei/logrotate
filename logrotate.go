package logrotate

import (
	"github.com/Niexiawei/logrotate/internal/ticker"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotateLog struct {
	file               *os.File
	logPath            string
	curLogLinkPath     string
	rotateTime         time.Duration
	maxAge             time.Duration
	deleteFileWildcard string
	timeTicker         *time.Ticker
	timeTickDate       ticker.TimeTickerDate
	mutex              *sync.Mutex
	rotate             <-chan time.Time // notify rotate event
	close              chan struct{}    // close file and write goroutine
}

func NewRotateLog(logPath string, opts ...Option) (*RotateLog, error) {
	rl := &RotateLog{
		mutex:        &sync.Mutex{},
		close:        make(chan struct{}, 1),
		logPath:      logPath,
		timeTickDate: ticker.TimeTickerDate{},
	}
	for _, opt := range opts {
		opt(rl)
	}

	rl.timeTicker = time.NewTicker(ticker.SetTime(rl.timeTickDate))

	if err := os.Mkdir(filepath.Dir(rl.logPath), 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if err := rl.rotateFile(time.Now()); err != nil {
		return nil, err
	}

	if rl.rotateTime != 0 {
		go rl.handleEvent()
	}

	return rl, nil
}

func (r *RotateLog) Write(b []byte) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	//print(r.file.Name(), string(b))
	n, err := r.file.Write(b)
	return n, err
}

func (r *RotateLog) Close() error {
	r.close <- struct{}{}
	return r.file.Close()
}

func (r *RotateLog) handleEvent() {
	for {
		select {
		case <-r.close:
			return
		case now := <-r.rotate:
			_ = r.rotateFile(now)
		case <-r.timeTicker.C:
			_ = r.rotateFile(time.Now())
		}
	}
}

func (r *RotateLog) rotateFile(now time.Time) error {
	if r.rotateTime != 0 {
		nextRotateTime := ticker.CalRotateTimeDuration(now, r.rotateTime)
		r.rotate = time.After(nextRotateTime)
	}

	latestLogPath := r.getLatestLogPath(now)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	file, err := os.OpenFile(latestLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if r.file != nil {
		_ = r.file.Close()
	}
	r.file = file

	if len(r.curLogLinkPath) > 0 {
		_ = os.Remove(r.curLogLinkPath)
		_ = os.Link(latestLogPath, r.curLogLinkPath)
	}

	if r.maxAge > 0 && len(r.deleteFileWildcard) > 0 { // at present
		go r.deleteExpiredFile(now)
	}

	return nil
}

// Judege expired by laste modify time
func (r *RotateLog) deleteExpiredFile(now time.Time) {
	cutoffTime := now.Add(-r.maxAge)
	matches, err := filepath.Glob(r.deleteFileWildcard)
	if err != nil {
		return
	}
	toUnlink := make([]string, 0, len(matches))
	for _, path := range matches {
		fileInfo, err := os.Stat(path)
		if err != nil {
			continue
		}

		if r.maxAge > 0 && fileInfo.ModTime().After(cutoffTime) {
			continue
		}

		if len(r.curLogLinkPath) > 0 && fileInfo.Name() == filepath.Base(r.curLogLinkPath) {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	for _, path := range toUnlink {
		_ = os.Remove(path)
	}
}

func (r *RotateLog) getLatestLogPath(t time.Time) string {
	return t.Format(r.logPath)
}
