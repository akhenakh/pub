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

func (self *Visibility) Scan(value interface{}) error {
	*self = Visibility(value.([]byte))
	return nil
}

func (self Visibility) Value() (driver.Value, error) {
	return string(self), nil
}
