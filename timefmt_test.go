package timefmt

import (
	"testing"
	"time"
)

func TestFormat(t *testing.T) {
	when := time.Date(2019, 5, 8, 17, 35, 18, 0, time.UTC)
	data := []struct {
		Format string
		Want   string
	}{
		{Format: "%Y", Want: "2019"},
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
	}
	for i, d := range data {
		got := Format(when, d.Format)
		if got != d.Want {
			t.Errorf("%d: date badly formatted (%s): want %s, got %s", i+1, d.Format, got, d.Want)
		}
	}
}
