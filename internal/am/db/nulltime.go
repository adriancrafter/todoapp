package db

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// NullTime struct
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Scan sql.Scanner interface
func (nt *NullTime) Scan(value interface{}) error {
	if value == nil {
		nt.Time, nt.Valid = time.Time{}, false
		return nil
	}
	nt.Valid = true

	switch v := value.(type) {
	case time.Time:
		nt.Time = v
	case []byte:
		t, err := time.Parse(time.RFC3339, string(v))
		if err != nil {
			return err
		}
		nt.Time = t
	case string:
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return err
		}
		nt.Time = t
	default:
		return fmt.Errorf("unsupported type for NullTime: %T", value)
	}
	return nil
}

// Value driver.Valuer interface
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
