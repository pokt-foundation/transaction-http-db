// Code generated by mockery v1.0.0. DO NOT EDIT.

package batch

import (
	context "context"

	types "github.com/pokt-foundation/transaction-db/types"
	mock "github.com/stretchr/testify/mock"
)

// MockRelayWriter is an autogenerated mock type for the RelayWriter type
type MockRelayWriter struct {
	mock.Mock
}

// WriteRelays provides a mock function with given fields: ctx, relays
func (_m *MockRelayWriter) WriteRelays(ctx context.Context, relays []types.Relay) error {
	ret := _m.Called(ctx, relays)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []types.Relay) error); ok {
		r0 = rf(ctx, relays)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
