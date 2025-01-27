package logrotate

import (
	"fmt"
	"github.com/Niexiawei/logrotate/internal/ticker"
	"path/filepath"
	"time"
)

type Option func(*RotateLog)

func WithRotateTime(duration time.Duration) Option {
	return func(r *RotateLog) {
		r.rotateTime = duration
	}
}

func WithCurLogLinkname(linkPath string) Option {
	return func(r *RotateLog) {
		r.curLogLinkPath = linkPath
	}
}

// WithDeleteExpiredFile Judege expired by laste modify time
// cutoffTime = now - maxAge
// Only delete satisfying file wildcard filename
func WithDeleteExpiredFile(maxAge time.Duration, fileWilCard string) Option {
	return func(r *RotateLog) {
		r.maxAge = maxAge
		r.deleteFileWildcard = fmt.Sprintf("%s%s%s", filepath.Dir(r.logPath), string([]byte{filepath.Separator}), fileWilCard)
	}
}

func WithTimeTickerDate(date ticker.TimeTickerDate) Option {
	return func(log *RotateLog) {
		log.timeTickDate = date
	}
}
