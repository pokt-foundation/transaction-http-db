package types

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

var (
	// this fields shpould be empty because they are set after db record is created
	shouldBeEmptySession = map[string]bool{
		"CreatedAt": true,
		"UpdatedAt": true,
	}

	ErrRepeatedSessionKey = errors.New("repeated session key")
)

type PocketSession struct {
	SessionKey       string    `json:"sessionKey"`
	SessionHeight    int       `json:"sessionHeight"`
	PortalRegionName string    `json:"portalRegionName"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

func (ps *PocketSession) Validate() error {
	structType := reflect.TypeOf(*ps)
	structVal := reflect.ValueOf(*ps)
	fieldNum := structVal.NumField()

	// fields are in the order they are declared on the struct
	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		isSet := field.IsValid() && !field.IsZero()

		if isSet {
			// shouldBeEmptyFields should never be set
			if shouldBeEmptySession[fieldName] {
				return fmt.Errorf("%s should not be set", fieldName)
			}
		}

		if !isSet {
			// shouldBeEmptyField can be empty
			// bools zero value is false which is a valid value
			if shouldBeEmptySession[fieldName] {
				continue
			}

			// if is not set and the field is none of the special cases it is an error
			return fmt.Errorf("%s is not set", fieldName)
		}
	}

	return nil
}
