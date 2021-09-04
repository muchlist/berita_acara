package dto

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type UppercaseString string

func (us UppercaseString) Value() (driver.Value, error) {
	return driver.Value(strings.ToUpper(string(us))), nil
}

func (us *UppercaseString) Scan(src interface{}) error {
	var source string
	switch src.(type) {
	case string:
		source = src.(string)
	case []byte:
		source = string(src.([]byte))
	default:
		return errors.New("Incompatible type for UppercaseString")
	}
	*us = UppercaseString(strings.ToUpper(source))
	return nil
}
