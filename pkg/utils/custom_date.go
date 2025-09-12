package utils

import (
	"strings"
	"time"
)

type CustomDate struct {
	time.Time
}

const dateLayout = "2006-01-02"

// UnmarshalJSON for parsing from JSON to struct
func (d *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// MarshalJSON for converting back to JSON format
func (d CustomDate) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte(`null`), nil
	}
	return []byte(`"` + d.Time.Format(dateLayout) + `"`), nil
}
