package timefmt

import (
	"time"
)

const (
	am = "AM"
	pm = "PM"
)

var timeNow = func() time.Time {
	return time.Now()
}

func isAM(t time.Time) bool {
	h, m, s := t.Hour(), t.Minute(), t.Second()
	if h == 0 && m == 0 && s == 0 {
		return true
	}
	return h < 12
}
