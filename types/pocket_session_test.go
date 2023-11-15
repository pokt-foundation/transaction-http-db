package types

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPocketSession_ValidateStruct(t *testing.T) {
	c := require.New(t)

	tests := []struct {
		name          string
		serviceRecord PocketSession
		err           error
	}{
		{
			name: "Success session",
			serviceRecord: PocketSession{
				SessionKey:       "key",
				SessionHeight:    1,
				PortalRegionName: "region",
			},
			err: nil,
		},
		{
			name: "failure session no key",
			serviceRecord: PocketSession{
				SessionKey:       "",
				SessionHeight:    1,
				PortalRegionName: "region",
			},
			err: errors.New("SessionKey is not set"),
		},
		{
			name: "failure session no height",
			serviceRecord: PocketSession{
				SessionKey:       "key",
				SessionHeight:    0,
				PortalRegionName: "region",
			},
			err: errors.New("SessionHeight is not set"),
		},
		{
			name: "failure session no region",
			serviceRecord: PocketSession{
				SessionKey:       "key",
				SessionHeight:    1,
				PortalRegionName: "",
			},
			err: errors.New("PortalRegionName is not set"),
		},
		{
			name: "failure session createdAt should not be set",
			serviceRecord: PocketSession{
				SessionKey:       "key",
				SessionHeight:    1,
				PortalRegionName: "region",
				CreatedAt:        time.Now(),
			},
			err: errors.New("CreatedAt should not be set"),
		},
	}

	for _, tt := range tests {
		c.Equal(tt.err, tt.serviceRecord.Validate())
	}
}
