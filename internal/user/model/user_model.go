package model

import (
	"database/sql/driver"
	"fmt"
)

type EnumHeightUnits string

const (
	EnumHeightUnitsCM   EnumHeightUnits = "CM"
	EnumHeightUnitsINCH EnumHeightUnits = "INCH"
)

func (e *EnumHeightUnits) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EnumHeightUnits(s)
	case string:
		*e = EnumHeightUnits(s)
	default:
		return fmt.Errorf("unsupported scan type for EnumHeightUnits: %T", src)
	}
	return nil
}

type NullEnumHeightUnits struct {
	EnumHeightUnits EnumHeightUnits
	Valid           bool // Valid is true if EnumHeightUnits is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEnumHeightUnits) Scan(value interface{}) error {
	if value == nil {
		ns.EnumHeightUnits, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EnumHeightUnits.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEnumHeightUnits) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EnumHeightUnits), nil
}

type EnumPreferences string

const (
	EnumPreferencesCARDIO EnumPreferences = "CARDIO"
	EnumPreferencesWEIGHT EnumPreferences = "WEIGHT"
)

func (e *EnumPreferences) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EnumPreferences(s)
	case string:
		*e = EnumPreferences(s)
	default:
		return fmt.Errorf("unsupported scan type for EnumPreferences: %T", src)
	}
	return nil
}

type NullEnumPreferences struct {
	EnumPreferences EnumPreferences
	Valid           bool // Valid is true if EnumPreferences is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEnumPreferences) Scan(value interface{}) error {
	if value == nil {
		ns.EnumPreferences, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EnumPreferences.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEnumPreferences) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EnumPreferences), nil
}

type EnumWeightUnits string

const (
	EnumWeightUnitsKG  EnumWeightUnits = "KG"
	EnumWeightUnitsLBS EnumWeightUnits = "LBS"
)

func (e *EnumWeightUnits) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = EnumWeightUnits(s)
	case string:
		*e = EnumWeightUnits(s)
	default:
		return fmt.Errorf("unsupported scan type for EnumWeightUnits: %T", src)
	}
	return nil
}

type NullEnumWeightUnits struct {
	EnumWeightUnits EnumWeightUnits
	Valid           bool // Valid is true if EnumWeightUnits is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullEnumWeightUnits) Scan(value interface{}) error {
	if value == nil {
		ns.EnumWeightUnits, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.EnumWeightUnits.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullEnumWeightUnits) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.EnumWeightUnits), nil
}

type User struct {
	ID             int
	Email          string
	HashedPassword string
	Username       *string
	UserImageUri   *string
	Weight         *int
	Height         *int
	WeightUnit     NullEnumWeightUnits
	HeightUnit     NullEnumHeightUnits
	Preference     NullEnumPreferences
}
