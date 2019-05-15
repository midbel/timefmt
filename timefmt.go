package timefmt

import (
	"fmt"
	"time"
)

const (
	am = "AM"
	pm = "PM"
)

type ParseError struct {
	Data    string
	Pattern string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("fail to parse %s as %s", e.Data, e.Pattern)
}

func parseError(data, pattern string) error {
	return ParseError{data, pattern}
}

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
