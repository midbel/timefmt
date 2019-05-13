package timefmt

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

func Format(w time.Time, pattern string) string {
	var buf bytes.Buffer

	r := strings.NewReader(pattern)
	for {
		b, err := r.ReadByte()
		if err != nil {
			break
		}
		if b == '%' {
			b, err := r.ReadByte()
			if err != nil {
				break
			}
			if b == 'E' || b == 'O' {
				b, _ = r.ReadByte() // ignore alternative conversion specifier
			}
			if fn, ok := formatters[b]; !ok || fn == nil {
				buf.WriteByte('%')
				buf.WriteByte(b)
			} else {
				fn(w, &buf)
			}
		} else {
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

var formatters map[byte]formatFunc

type formatFunc func(time.Time, io.Writer)

func init() {
	formatters = map[byte]formatFunc{
		'a': formatDayName(true), //dayShortname,
		'A': formatDayName(false),
		'b': formatMonthName(true),
		'B': formatMonthName(false),
		'c': formatDatetime,
		'C': formatCentury,
		'd': formatDay("%02d"),
		'D': formatDateAlt,
		'e': formatDay("%2d"),
		'F': formatDate,
		'g': nil,
		'G': formatYearISO,
		'h': formatMonthName(true),
		'H': formatHour(24, "%02d"),
		'I': formatHour(12, "%02d"),
		'j': formatYearDay,
		'k': formatHour(24, "%2d"),
		'l': formatHour(12, "%2d"),
		'm': formatMonth,
		'M': formatMinute,
		'p': formatAPM(false),
		'P': formatAPM(true),
		'r': formatIMSP,
		'R': formatShortTime,
		's': formatTimestamp,
		'S': formatSeconds,
		'T': formatLongTime,
		'u': formatWeekday(false),
		'V': formatWeekNumberISO,
		'w': formatWeekday(true),
		'x': formatDate,
		'X': formatLongTime,
		'y': formatShortYear,
		'Y': formatLongYear,
		'n': formatChar("\n"),
		't': formatChar("\t"),
		'%': formatChar("%"),
	}
}

func formatChar(c string) formatFunc {
	return func(_ time.Time, w io.Writer) {
		io.WriteString(w, c)
	}
}

func formatWeekday(sunday bool) formatFunc {
	return func(t time.Time, w io.Writer) {
		wk := t.Weekday()
		if sunday {
			wk--
		}
		io.WriteString(w, fmt.Sprint(int(wk)))
	}
}

func formatWeekNumberISO(t time.Time, w io.Writer) {
	_, wk := t.ISOWeek()
	io.WriteString(w, fmt.Sprintf("%02d", wk))
}

func formatYearISO(t time.Time, w io.Writer) {
	yk, _ := t.ISOWeek()
	io.WriteString(w, fmt.Sprintf("%04d", yk))
}

func formatIMSP(t time.Time, w io.Writer) {
	formatters['I'](t, w)
	io.WriteString(w, ":")
	formatMinute(t, w)
	io.WriteString(w, ":")
	formatSeconds(t, w)
	io.WriteString(w, " ")
	formatters['p'](t, w)
}

func formatTimestamp(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprint(t.Unix()))
}

func formatShortTime(t time.Time, w io.Writer) {
	formatters['H'](t, w)
	io.WriteString(w, ":")
	formatMinute(t, w)
}

func formatLongTime(t time.Time, w io.Writer) {
	formatShortTime(t, w)
	io.WriteString(w, ":")
	formatSeconds(t, w)
}

func formatShortTimeAlt(t time.Time, w io.Writer) {
	formatters['I'](t, w)
	io.WriteString(w, ":")
	formatMinute(t, w)
}

func formatLongTimeAlt(t time.Time, w io.Writer) {
	formatShortTimeAlt(t, w)
	io.WriteString(w, ":")
	formatSeconds(t, w)
}

func formatDateAlt(t time.Time, w io.Writer) {
	formatters['d'](t, w)
	io.WriteString(w, "/")
	formatMonth(t, w)
	io.WriteString(w, "/")
	formatLongYear(t, w)
}

func formatDate(t time.Time, w io.Writer) {
	formatLongYear(t, w)
	io.WriteString(w, "-")
	formatMonth(t, w)
	io.WriteString(w, "-")
	formatters['d'](t, w)
}

func formatShortYear(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Year()%100))
}

func formatLongYear(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%04d", t.Year()))
}

func formatSeconds(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Second()))
}

func formatMinute(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Minute()))
}

func formatMonth(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Month()))
}

func formatYearDay(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%03d", t.YearDay()))
}

func formatHour(clock int, pattern string) formatFunc {
	return func(t time.Time, w io.Writer) {
		h := t.Hour()
		if h >= clock {
			h -= clock
		}
		io.WriteString(w, fmt.Sprintf(pattern, h))
	}
}

func formatDay(pattern string) formatFunc {
	return func(t time.Time, w io.Writer) {
		io.WriteString(w, fmt.Sprintf(pattern, t.Day()))
	}
}

func formatCentury(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Year()/100))
}

func formatDayName(abbr bool) formatFunc {
	return func(t time.Time, w io.Writer) {
		n := t.Weekday().String()
		if abbr {
			n = n[:3]
		}
		io.WriteString(w, n)
	}
}

func formatMonthName(abbr bool) formatFunc {
	return func(t time.Time, w io.Writer) {
		n := t.Month().String()
		if abbr {
			n = n[:3]
		}
		io.WriteString(w, n)
	}
}

func formatDatetime(t time.Time, w io.Writer) {
	formatters['a'](t, w)
	io.WriteString(w, " ")
	formatters['d'](t, w)
	io.WriteString(w, " ")
	formatters['b'](t, w)
	io.WriteString(w, " ")
	formatLongYear(t, w)
	io.WriteString(w, " ")
	formatLongTimeAlt(t, w)
	io.WriteString(w, " ")
	formatters['p'](t, w)
}

func formatAPM(lower bool) formatFunc {
	return func(t time.Time, w io.Writer) {
		str := pm
		if isAM(t) {
			str = am
		}
		if lower {
			str = strings.ToLower(str)
		}
		io.WriteString(w, str)
	}
}
