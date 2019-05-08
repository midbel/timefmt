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
			fn, ok := tokens[b]
			if !ok || fn == nil {
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

const (
	am = "AM"
	pm = "PM"
)

type formatFunc func(time.Time, io.Writer)

var tokens = map[byte]formatFunc{
	'a': dayShortname,
	'A': dayLongname,
	'b': monthShortname,
	'B': monthLongname,
	'c': formatDatetime,
	'C': formatCentury,
	'd': formatDay,
	'D': nil,
	'e': nil,
	'F': formatDate,
	'g': nil,
	'G': formatYearISO,
	'H': formatHour24,
	'I': formatHour12,
	'j': formatYearDay,
	'm': formatMonth,
	'M': formatMinute,
	'p': formatUpperAMPM,
	'P': formatLowerAMPM,
	'r': formatIMSP,
	'R': formatShortTime,
	's': formatTimestamp,
	'S': formatSeconds,
	'T': formatLongTime,
	'V': formatWeekNumberISO,
	'x': formatDate,
	'X': formatLongTime,
	'y': formatShortYear,
	'Y': formatLongYear,
	'n': func(_ time.Time, w io.Writer) { io.WriteString(w, "\n") },
	't': func(_ time.Time, w io.Writer) { io.WriteString(w, "\t") },
	'%': func(_ time.Time, w io.Writer) { io.WriteString(w, "%") },
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
	formatHour12(t, w)
	io.WriteString(w, ":")
	formatMinute(t, w)
	io.WriteString(w, ":")
	formatSeconds(t, w)
	io.WriteString(w, " ")
	formatUpperAMPM(t, w)
}

func formatTimestamp(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprint(t.Unix()))
}

func formatShortTime(t time.Time, w io.Writer) {
	formatHour24(t, w)
	io.WriteString(w, ":")
	formatMinute(t, w)
}

func formatLongTime(t time.Time, w io.Writer) {
	formatShortTime(t, w)
	io.WriteString(w, ":")
	formatSeconds(t, w)
}

func formatDate(t time.Time, w io.Writer) {
	formatLongYear(t, w)
	io.WriteString(w, "-")
	formatMonth(t, w)
	io.WriteString(w, "-")
	formatDay(t, w)
}

func formatShortYear(t time.Time, w io.Writer) {

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

func formatHour24(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Hour()))
}

func formatHour12(t time.Time, w io.Writer) {
	h := t.Hour()
	if h >= 12 {
		h -= 12
	}
	io.WriteString(w, fmt.Sprintf("%02d", h))
}

func formatDay(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Day()))
}

func formatCentury(t time.Time, w io.Writer) {
	io.WriteString(w, fmt.Sprintf("%02d", t.Year()/100))
}

func dayShortname(t time.Time, w io.Writer) {
	wk := t.Weekday().String()
	io.WriteString(w, wk[:3])
}

func dayLongname(t time.Time, w io.Writer) {
	wk := t.Weekday()
	io.WriteString(w, wk.String())
}

func monthLongname(t time.Time, w io.Writer) {
	mn := t.Month()
	io.WriteString(w, mn.String())
}

func monthShortname(t time.Time, w io.Writer) {
	mn := t.Month().String()
	io.WriteString(w, mn[:3])
}

func formatDatetime(t time.Time, w io.Writer) {
	formatDate(t, w)
	io.WriteString(w, " ")
	formatLongTime(t, w)
}

func formatLowerAMPM(t time.Time, w io.Writer) {
	str := strings.ToLower(pm)
	if isAM(t) {
		str = strings.ToLower(am)
	}
	io.WriteString(w, str)
}

func formatUpperAMPM(t time.Time, w io.Writer) {
	str := pm
	if isAM(t) {
		str = am
	}
	io.WriteString(w, str)
}

func isAM(t time.Time) bool {
	h, m, s := t.Hour(), t.Minute(), t.Second()
	if h == 0 && m == 0 && s == 0 {
		return true
	}
	return h < 12
}
