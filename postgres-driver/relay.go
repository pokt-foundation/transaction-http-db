package postgresdriver

import (
	"context"
	"strings"
	"time"

	"github.com/pokt-foundation/transaction-http-db/types"
)

const chainMethodIDSeparator = ","

func (d *PostgresDriver) WriteRelay(ctx context.Context, relay types.Relay) error {
	createdAt := time.Now()

	_, err := d.InsertRelays(ctx, []InsertRelaysParams{
		{
			PoktChainID:              relay.PoktChainID,
			EndpointID:               relay.EndpointID,
			SessionKey:               relay.SessionKey,
			ProtocolAppPublicKey:     relay.ProtocolAppPublicKey,
			RelaySourceUrl:           newText(relay.RelaySourceURL),
			PoktNodeAddress:          newText(relay.PoktNodeAddress),
			PoktNodeDomain:           newText(relay.PoktNodeDomain),
			PoktNodePublicKey:        newText(relay.PoktNodePublicKey),
			RelayStartDatetime:       newTimestamp(relay.RelayStartDatetime),
			RelayReturnDatetime:      newTimestamp(relay.RelayReturnDatetime),
			IsError:                  relay.IsError,
			ErrorCode:                newInt4(int32(relay.ErrorCode), false),
			ErrorName:                newText(relay.ErrorName),
			ErrorMessage:             newText(relay.ErrorMessage),
			ErrorType:                newText(relay.ErrorType),
			ErrorSource:              newNullErrorSourcesEnum(ErrorSourcesEnum(relay.ErrorSource)),
			RelayRoundtripTime:       relay.RelayRoundtripTime,
			RelayChainMethodIds:      strings.Join(relay.RelayChainMethodIDs, chainMethodIDSeparator),
			RelayDataSize:            int32(relay.RelayDataSize),
			RelayPortalTripTime:      relay.RelayPortalTripTime,
			RelayNodeTripTime:        relay.RelayNodeTripTime,
			RelayUrlIsPublicEndpoint: relay.RelayURLIsPublicEndpoint,
			PortalRegionName:         relay.PortalRegionName,
			IsAltruistRelay:          relay.IsAltruistRelay,
			RequestID:                relay.RequestID,
			PoktTxID:                 newText(relay.PoktTxID),
			IsUserRelay:              relay.IsUserRelay,
			GigastakeAppID:           newText(relay.GigastakeAppID),
			CreatedAt:                newTimestamp(createdAt),
			UpdatedAt:                newTimestamp(createdAt),
			BlockingPlugin:           newText(relay.BlockingPlugin),
		},
	})

	return err
}

func (d *PostgresDriver) WriteRelays(ctx context.Context, relays []*types.Relay) error {
	createdAt := time.Now()

	relayParams := make([]InsertRelaysParams, 0, len(relays))

	for _, relay := range relays {
		relayParams = append(relayParams, InsertRelaysParams{
			PoktChainID:              relay.PoktChainID,
			EndpointID:               relay.EndpointID,
			SessionKey:               relay.SessionKey,
			ProtocolAppPublicKey:     relay.ProtocolAppPublicKey,
			RelaySourceUrl:           newText(relay.RelaySourceURL),
			PoktNodeAddress:          newText(relay.PoktNodeAddress),
			PoktNodeDomain:           newText(relay.PoktNodeDomain),
			PoktNodePublicKey:        newText(relay.PoktNodePublicKey),
			RelayStartDatetime:       newTimestamp(relay.RelayStartDatetime),
			RelayReturnDatetime:      newTimestamp(relay.RelayReturnDatetime),
			IsError:                  relay.IsError,
			ErrorCode:                newInt4(int32(relay.ErrorCode), false),
			ErrorName:                newText(relay.ErrorName),
			ErrorMessage:             newText(relay.ErrorMessage),
			ErrorType:                newText(relay.ErrorType),
			ErrorSource:              newNullErrorSourcesEnum(ErrorSourcesEnum(relay.ErrorSource)),
			RelayRoundtripTime:       relay.RelayRoundtripTime,
			RelayChainMethodIds:      strings.Join(relay.RelayChainMethodIDs, chainMethodIDSeparator),
			RelayDataSize:            int32(relay.RelayDataSize),
			RelayPortalTripTime:      relay.RelayPortalTripTime,
			RelayNodeTripTime:        relay.RelayNodeTripTime,
			RelayUrlIsPublicEndpoint: relay.RelayURLIsPublicEndpoint,
			PortalRegionName:         relay.PortalRegionName,
			IsAltruistRelay:          relay.IsAltruistRelay,
			RequestID:                relay.RequestID,
			PoktTxID:                 newText(relay.PoktTxID),
			IsUserRelay:              relay.IsUserRelay,
			GigastakeAppID:           newText(relay.GigastakeAppID),
			CreatedAt:                newTimestamp(createdAt),
			UpdatedAt:                newTimestamp(createdAt),
			BlockingPlugin:           newText(relay.BlockingPlugin),
		})
	}

	_, err := d.InsertRelays(ctx, relayParams)

	return err
}

func (d *PostgresDriver) ReadRelay(ctx context.Context, relayID int) (types.Relay, error) {
	relay, err := d.SelectRelay(ctx, int64(relayID))
	if err != nil {
		return types.Relay{}, err
	}

	return types.Relay{
		RelayID:                  int(relay.ID),
		PoktChainID:              relay.PoktChainID,
		EndpointID:               relay.EndpointID,
		SessionKey:               relay.SessionKey,
		ProtocolAppPublicKey:     relay.ProtocolAppPublicKey,
		RelaySourceURL:           relay.RelaySourceUrl.String,
		PoktNodeAddress:          relay.PoktNodeAddress.String,
		PoktNodeDomain:           relay.PoktNodeDomain.String,
		PoktNodePublicKey:        relay.PoktNodePublicKey.String,
		RelayStartDatetime:       relay.RelayStartDatetime.Time,
		RelayReturnDatetime:      relay.RelayReturnDatetime.Time,
		IsError:                  relay.IsError,
		ErrorCode:                int(relay.ErrorCode.Int32),
		ErrorName:                relay.ErrorName.String,
		ErrorMessage:             relay.ErrorMessage.String,
		ErrorType:                relay.ErrorType.String,
		ErrorSource:              types.ErrorSource(relay.ErrorSource.ErrorSourcesEnum),
		RelayRoundtripTime:       relay.RelayRoundtripTime,
		RelayChainMethodIDs:      strings.Split(relay.RelayChainMethodIds, ","),
		RelayDataSize:            int(relay.RelayDataSize),
		RelayPortalTripTime:      relay.RelayPortalTripTime,
		RelayNodeTripTime:        relay.RelayNodeTripTime,
		RelayURLIsPublicEndpoint: relay.RelayUrlIsPublicEndpoint,
		PortalRegionName:         relay.PortalRegionName,
		IsAltruistRelay:          relay.IsAltruistRelay,
		RequestID:                relay.RequestID,
		IsUserRelay:              relay.IsUserRelay,
		PoktTxID:                 relay.PoktTxID.String,
		GigastakeAppID:           relay.GigastakeAppID.String,
		CreatedAt:                relay.CreatedAt.Time,
		UpdatedAt:                relay.UpdatedAt.Time,
		Session: types.PocketSession{
			SessionKey:    relay.SessionKey,
			SessionHeight: int(relay.SessionHeight),
			CreatedAt:     relay.CreatedAt_2.Time,
			UpdatedAt:     relay.UpdatedAt_2.Time,
		},
		Region: types.PortalRegion{
			PortalRegionName: relay.PortalRegionName,
		},
	}, nil
}
