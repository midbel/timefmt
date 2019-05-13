package timefmt

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2019, 5, 8, 17, 35, 18, 0, time.Local)
	}
	data := []struct {
		Format string
		Data   string
		Want   time.Time
	}{
		{Format: "%Y", Data: "2019", Want: time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local)},
		{Format: "%Y-%m-%d", Data: "2019-05-08", Want: time.Date(2019, 5, 8, 0, 0, 0, 0, time.Local)},
		{Format: "%d/%m/%Y", Data: "08/05/2019", Want: time.Date(2019, 5, 8, 0, 0, 0, 0, time.Local)},
		{Format: "%D", Data: "08/05/2019", Want: time.Date(2019, 5, 8, 0, 0, 0, 0, time.Local)},
		{Format: "%H:%M", Data: "17:35", Want: timeNow().Truncate(60 * time.Second)},
		{Format: "%R", Data: "17:35", Want: timeNow().Truncate(60 * time.Second)},
		{Format: "%H:%M:%S", Data: "17:35:18", Want: timeNow()},
		{Format: "%T", Data: "17:35:18", Want: timeNow()},
	}
	for i, d := range data {
		got := Parse(d.Data, d.Format)
		if !got.Equal(d.Want) {
			// t.Log(got, d.Want)
			t.Errorf("%d: date badly parsed: want %s, got %s", i+1, Format(d.Want, d.Format), Format(got, d.Format))
		}
	}
}

func TestFormat(t *testing.T) {
	timeNow = func() time.Time {
		return time.Date(2019, 5, 8, 17, 35, 18, 0, time.Local)
	}
	when := timeNow()
	data := []struct {
		Format string
		Want   string
	}{
		{Format: "%Y", Want: "2019"},
		{Format: "%y", Want: "19"},
		{Format: "%G", Want: "2019"},
		{Format: "%Y/%j", Want: "2019/128"},
		{Format: "%C", Want: "20"},
		{Format: "%d", Want: "08"},
		{Format: "%e", Want: " 8"},
		{Format: "%F", Want: "2019-05-08"},
		{Format: "%Y-%m-%d", Want: "2019-05-08"},
		{Format: "%H:%M", Want: "17:35"},
		{Format: "%R", Want: "17:35"},
		{Format: "%H:%M:%S", Want: "17:35:18"},
		{Format: "%T", Want: "17:35:18"},
		{Format: "%D", Want: "08/05/2019"},
		{Format: "%a(%u)", Want: "Wed(3)"},
		{Format: "%A(%w)", Want: "Wednesday(2)"},
		{Format: "%b", Want: "May"},
		{Format: "%B", Want: "May"},
		{Format: "%c", Want: "Wed 08 May 2019 05:35:18 PM"},
	}
	for i, d := range data {
		got := Format(when, d.Format)
		if got != d.Want {
			t.Errorf("%d: date badly formatted (%s): want %s, got %s", i+1, d.Format, got, d.Want)
		}
	}
}
