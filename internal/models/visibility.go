package models

import (
	"database/sql/driver"
	"errors"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Visibility string

const (
	PublicVisibility   Visibility = "public"
	UnlistedVisibility Visibility = "unlisted"
	PrivateVisibility  Visibility = "private"
	DirectVisibility   Visibility = "direct"
	LimitedVisibility  Visibility = "limited"
)

func (v *Visibility) Scan(value interface{}) error {
	var pv Visibility
	if value == nil {
		*v = ""
		return nil
	}
	st, ok := value.([]uint8)
	if !ok {
		return errors.New("Invalid data for visibility")
	}

	pv = Visibility(string(st))

	switch pv {
	case PublicVisibility, UnlistedVisibility, PrivateVisibility, DirectVisibility, LimitedVisibility:
		*v = pv
		return nil
	}
	return nil
}

func (v Visibility) Value() (driver.Value, error) {
	switch v {
	case PublicVisibility, UnlistedVisibility, PrivateVisibility, DirectVisibility, LimitedVisibility:
		return string(v), nil
	}
	return nil, errors.New("Invalid visiblity value")
}
