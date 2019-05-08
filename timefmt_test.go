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
		{Format: "%Y/%j", Want: "2019/128"},
	}
	for i, d := range data {
		got := Format(when, d.Format)
		if got != d.Want {
			t.Errorf("%d: date badly formatted (%s): want %s, got %s", i+1, d.Format, got, d.Want)
		}
	}
}
