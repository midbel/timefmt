package timefmt

import (
	"strconv"
	"strings"
	"time"
)

func Parse(str, pattern string) time.Time {
	r := strings.NewReader(pattern)

	var (
		dt  datetime
		pos int
	)
	ds := []byte(str)
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if b == '%' {
			if b, err = r.ReadByte(); err != nil {
				break
			}
			if b == 'E' || b == 'O' {
				b, _ = r.ReadByte() // ignore alternative conversion specifier
			}
			if fn, ok := converters[b]; ok {
				n := fn(&dt, ds[pos:])
				if n <= 0 {
					break
				}
				pos += n
			}
		} else {
			if b != ds[pos] {
				break
			}
			pos++
		}
	}
	return dt.Time()
}

type datetime struct {
	sec  int // Seconds (0-60)
	min  int // Minutes (0-59)
	hour int // Hours (0-23)
	day  int // Day of the month (1-31)
	mon  int // Month (0-11)
	year int // Year - 1900

	wday  int // Day of the week (0-6, Sunday = 0)
	yday  int // Day in the year (0-365, 1 Jan = 0)
	isdst int // Daylight saving time
}

func (d datetime) Time() time.Time {
	if d.year == 0 && d.mon == 0 && d.day == 0 {
		year, month, day := timeNow().Date()
		d.year, d.mon, d.day = year, int(month), day
	}
	if d.mon == 0 {
		d.mon = 1
	}
	if d.day == 0 {
		d.day = 1
	}
	t := time.Date(d.year, time.Month(d.mon), d.day, d.hour, d.min, d.sec, 0, time.Local)
	return t
}

type parseFunc func(*datetime, []byte) int

var converters map[byte]parseFunc

func init() {
	converters = map[byte]parseFunc{
		'd': parseDay,
		'D': parseDateAlt,
		'H': parseHour,
		'm': parseMonth,
		'M': parseMinute,
		'S': parseSecond,
		'R': parseShortTime,
		'T': parseLongTime,
		'Y': parseLongYear,
	}
}

func parseHour(dt *datetime, bs []byte) int {
	dt.hour = parseInt(bs, 2)
	return 2
}

func parseMinute(dt *datetime, bs []byte) int {
	dt.min = parseInt(bs, 2)
	return 2
}

func parseSecond(dt *datetime, bs []byte) int {
	dt.sec = parseInt(bs, 2)
	return 2
}

func parseLongTime(dt *datetime, bs []byte) int {
	pos := parseShortTime(dt, bs)
	if bs[pos] != ':' {
		return 0
	}
	pos++
	pos += parseSecond(dt, bs[pos:])
	return pos
}

func parseShortTime(dt *datetime, bs []byte) int {
	var pos int
	pos += parseHour(dt, bs)
	if bs[pos] != ':' {
		return 0
	}
	pos++

	pos += parseMinute(dt, bs[pos:])
	return pos
}

func parseDay(dt *datetime, bs []byte) int {
	dt.day = parseInt(bs, 2)
	return 2
}

func parseMonth(dt *datetime, bs []byte) int {
	dt.mon = parseInt(bs, 2)
	return 2
}

func parseLongYear(dt *datetime, bs []byte) int {
	dt.year = parseInt(bs, 4)
	return 4
}

func parseDateAlt(dt *datetime, bs []byte) int {
	var pos int
	pos += parseDay(dt, bs)
	if bs[pos] != '/' {
		return 0
	}
	pos++

	pos += parseMonth(dt, bs[pos:])
	if bs[pos] != '/' {
		return 0
	}
	pos++

	pos += parseLongYear(dt, bs[pos:])
	return pos
}

func parseInt(bs []byte, n int) int {
	v, err := strconv.ParseInt(string(bs[:n]), 10, 64)
	if err != nil {
		v = -1
	}
	return int(v)
}
