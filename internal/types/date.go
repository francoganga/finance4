package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Date struct {
	time.Time
}

func NewDate(t time.Time) *Date {
	return &Date{
		Time: t,
	}
}

func (d Date) Value() (driver.Value, error) {

	return d.Time.Format("2006-01-02"), nil
}

func (d *Date) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("Expected string, got null")
	}

	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("Value is not a string")
	}

	tv, err := time.Parse("2006-01-02", v)
	if err != nil {
		return fmt.Errorf("Scan Error: Could not parse %s as a date", v)
	}

	d.Time = tv

	return nil
}

