// Code generated by mockery v1.0.0. DO NOT EDIT.

package batch

import (
	context "context"

	types "github.com/pokt-foundation/transaction-db/types"
	mock "github.com/stretchr/testify/mock"
)

// MockServiceRecordWriter is an autogenerated mock type for the ServiceRecordWriter type
type MockServiceRecordWriter struct {
	mock.Mock
}

// WriteServiceRecords provides a mock function with given fields: ctx, serviceRecords
func (_m *MockServiceRecordWriter) WriteServiceRecords(ctx context.Context, serviceRecords []types.ServiceRecord) error {
	ret := _m.Called(ctx, serviceRecords)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []types.ServiceRecord) error); ok {
		r0 = rf(ctx, serviceRecords)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
