package postgresdriver

import (
	"context"
	"strings"
	"time"

	"github.com/pokt-foundation/transaction-http-db/types"
)

const (
	errMessageDuplicateSessionKey = `duplicate key value violates unique constraint "pocket_session_session_key_key"`
)

func (d *PostgresDriver) WriteSession(ctx context.Context, session types.PocketSession) error {
	now := time.Now()

	err := d.InsertPocketSession(ctx, InsertPocketSessionParams{
		SessionKey:       session.SessionKey,
		SessionHeight:    int32(session.SessionHeight),
		PortalRegionName: session.PortalRegionName,
		CreatedAt:        newTimestamp(now),
		UpdatedAt:        newTimestamp(now),
	})
	if err != nil {
		if strings.Contains(err.Error(), errMessageDuplicateSessionKey) {
			return types.ErrRepeatedSessionKey
		}

		return err
	}

	return nil
}
