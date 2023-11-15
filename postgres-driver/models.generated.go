// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package postgresdriver

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type ErrorSourcesEnum string

const (
	ErrorSourcesEnumInternal ErrorSourcesEnum = "internal"
	ErrorSourcesEnumExternal ErrorSourcesEnum = "external"
)

func (e *ErrorSourcesEnum) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = ErrorSourcesEnum(s)
	case string:
		*e = ErrorSourcesEnum(s)
	default:
		return fmt.Errorf("unsupported scan type for ErrorSourcesEnum: %T", src)
	}
	return nil
}

type NullErrorSourcesEnum struct {
	ErrorSourcesEnum ErrorSourcesEnum `json:"errorSourcesEnum"`
	Valid            bool             `json:"valid"` // Valid is true if ErrorSourcesEnum is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullErrorSourcesEnum) Scan(value interface{}) error {
	if value == nil {
		ns.ErrorSourcesEnum, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.ErrorSourcesEnum.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullErrorSourcesEnum) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.ErrorSourcesEnum), nil
}

type PocketSession struct {
	ID               int64            `json:"id"`
	SessionKey       string           `json:"sessionKey"`
	SessionHeight    int32            `json:"sessionHeight"`
	PortalRegionName string           `json:"portalRegionName"`
	CreatedAt        pgtype.Timestamp `json:"createdAt"`
	UpdatedAt        pgtype.Timestamp `json:"updatedAt"`
}

type PortalRegion struct {
	PortalRegionName string `json:"portalRegionName"`
}

type Relay struct {
	ID                       int64                `json:"id"`
	PoktChainID              string               `json:"poktChainId"`
	EndpointID               string               `json:"endpointId"`
	SessionKey               string               `json:"sessionKey"`
	ProtocolAppPublicKey     string               `json:"protocolAppPublicKey"`
	RelaySourceUrl           pgtype.Text          `json:"relaySourceUrl"`
	PoktNodeAddress          pgtype.Text          `json:"poktNodeAddress"`
	PoktNodeDomain           pgtype.Text          `json:"poktNodeDomain"`
	PoktNodePublicKey        pgtype.Text          `json:"poktNodePublicKey"`
	RelayStartDatetime       pgtype.Timestamp     `json:"relayStartDatetime"`
	RelayReturnDatetime      pgtype.Timestamp     `json:"relayReturnDatetime"`
	IsError                  bool                 `json:"isError"`
	ErrorCode                pgtype.Int4          `json:"errorCode"`
	ErrorName                pgtype.Text          `json:"errorName"`
	ErrorMessage             pgtype.Text          `json:"errorMessage"`
	ErrorSource              NullErrorSourcesEnum `json:"errorSource"`
	ErrorType                pgtype.Text          `json:"errorType"`
	RelayRoundtripTime       float64              `json:"relayRoundtripTime"`
	RelayChainMethodIds      string               `json:"relayChainMethodIds"`
	RelayDataSize            int32                `json:"relayDataSize"`
	RelayPortalTripTime      float64              `json:"relayPortalTripTime"`
	RelayNodeTripTime        float64              `json:"relayNodeTripTime"`
	RelayUrlIsPublicEndpoint bool                 `json:"relayUrlIsPublicEndpoint"`
	PortalRegionName         string               `json:"portalRegionName"`
	IsAltruistRelay          bool                 `json:"isAltruistRelay"`
	IsUserRelay              bool                 `json:"isUserRelay"`
	RequestID                string               `json:"requestId"`
	PoktTxID                 pgtype.Text          `json:"poktTxId"`
	GigastakeAppID           pgtype.Text          `json:"gigastakeAppId"`
	CreatedAt                pgtype.Timestamp     `json:"createdAt"`
	UpdatedAt                pgtype.Timestamp     `json:"updatedAt"`
	BlockingPlugin           pgtype.Text          `json:"blockingPlugin"`
}

type ServiceRecord struct {
	ID                     int64            `json:"id"`
	NodePublicKey          string           `json:"nodePublicKey"`
	PoktChainID            string           `json:"poktChainId"`
	SessionKey             string           `json:"sessionKey"`
	RequestID              string           `json:"requestId"`
	PortalRegionName       string           `json:"portalRegionName"`
	Latency                float64          `json:"latency"`
	Tickets                int32            `json:"tickets"`
	Result                 string           `json:"result"`
	Available              bool             `json:"available"`
	Successes              int32            `json:"successes"`
	Failures               int32            `json:"failures"`
	P90SuccessLatency      float64          `json:"p90SuccessLatency"`
	MedianSuccessLatency   float64          `json:"medianSuccessLatency"`
	WeightedSuccessLatency float64          `json:"weightedSuccessLatency"`
	SuccessRate            float64          `json:"successRate"`
	CreatedAt              pgtype.Timestamp `json:"createdAt"`
	UpdatedAt              pgtype.Timestamp `json:"updatedAt"`
}
