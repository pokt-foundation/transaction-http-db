package types

import (
	"fmt"
	"reflect"
	"time"
)

var (
	// this fields shpould be empty because they are set after db record is created
	shouldBeEmptyServiceRecordField = map[string]bool{
		"ServiceRecordID": true,
		"CreatedAt":       true,
		"UpdatedAt":       true,
	}

	serviceRecordOptionalFields = map[string]bool{
		"Latency":   true,
		"Result":    true,
		"Successes": true,
		"Failures":  true,
	}
)

type ServiceRecord struct {
	ServiceRecordID        int       `json:"serviceRecordID"`
	NodePublicKey          string    `json:"nodePublicKey"`
	PoktChainID            string    `json:"poktChainID"`
	SessionKey             string    `json:"sessionKey"`
	RequestID              string    `json:"requestID"`
	PortalRegionName       string    `json:"portalRegionName"`
	Latency                float64   `json:"latency"`
	Tickets                int       `json:"tickets"`
	Result                 string    `json:"result"`
	Available              bool      `json:"available"`
	Successes              int       `json:"successes"`
	Failures               int       `json:"failures"`
	P90SuccessLatency      float64   `json:"p90SuccessLatency"`
	MedianSuccessLatency   float64   `json:"medianSuccessLatency"`
	WeightedSuccessLatency float64   `json:"weightedSuccessLatency"`
	SuccessRate            float64   `json:"successRate"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

func (sr *ServiceRecord) Validate() (err error) {
	structType := reflect.TypeOf(*sr)
	structVal := reflect.ValueOf(*sr)
	fieldNum := structVal.NumField()

	// fields are in the order they are declared on the struct
	for i := 0; i < fieldNum; i++ {
		field := structVal.Field(i)
		fieldName := structType.Field(i).Name

		isSet := field.IsValid() && !field.IsZero()

		if isSet {
			// shouldBeEmptyFields should never be set
			if shouldBeEmptyServiceRecordField[fieldName] {
				return fmt.Errorf("%s should not be set", fieldName)
			}
		}

		if !isSet {
			// shouldBeEmptyField can be empty
			// bools zero value is false which is a valid value
			if shouldBeEmptyServiceRecordField[fieldName] ||
				field.Kind() == reflect.Bool ||
				serviceRecordOptionalFields[fieldName] {
				continue
			}

			// if is not set and the field is none of the special cases it is an error
			return fmt.Errorf("%s is not set", fieldName)
		}
	}

	return nil
}
