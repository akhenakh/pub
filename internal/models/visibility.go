package models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type (
	Visibility    string
	Conversations struct {
		db *gorm.DB
	}
)

const (
	PublicVisibility   Visibility = "public"
	UnlistedVisibility Visibility = "unlisted"
	PrivateVisibility  Visibility = "private"
	DirectVisibility   Visibility = "direct"
	LimitedVisibility  Visibility = "limited"
)

func (v *Visibility) Scan(value interface{}) error {
	*v = Visibility(value.([]byte))
	return nil
}

func (v Visibility) Value() (driver.Value, error) {
	return string(v), nil
}
