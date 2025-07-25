// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"context"

	"github.com/LiquidCats/rater/internal/adapter/repository/database/postgres"
	mock "github.com/stretchr/testify/mock"
)

// NewPairDatabase creates a new instance of PairDatabase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPairDatabase(t interface {
	mock.TestingT
	Cleanup(func())
}) *PairDatabase {
	mock := &PairDatabase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// PairDatabase is an autogenerated mock type for the PairDatabase type
type PairDatabase struct {
	mock.Mock
}

type PairDatabase_Expecter struct {
	mock *mock.Mock
}

func (_m *PairDatabase) EXPECT() *PairDatabase_Expecter {
	return &PairDatabase_Expecter{mock: &_m.Mock}
}

// CountPairs provides a mock function for the type PairDatabase
func (_mock *PairDatabase) CountPairs(ctx context.Context) (int64, error) {
	ret := _mock.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CountPairs")
	}

	var r0 int64
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context) (int64, error)); ok {
		return returnFunc(ctx)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context) int64); ok {
		r0 = returnFunc(ctx)
	} else {
		r0 = ret.Get(0).(int64)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = returnFunc(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// PairDatabase_CountPairs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CountPairs'
type PairDatabase_CountPairs_Call struct {
	*mock.Call
}

// CountPairs is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PairDatabase_Expecter) CountPairs(ctx interface{}) *PairDatabase_CountPairs_Call {
	return &PairDatabase_CountPairs_Call{Call: _e.mock.On("CountPairs", ctx)}
}

func (_c *PairDatabase_CountPairs_Call) Run(run func(ctx context.Context)) *PairDatabase_CountPairs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *PairDatabase_CountPairs_Call) Return(n int64, err error) *PairDatabase_CountPairs_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *PairDatabase_CountPairs_Call) RunAndReturn(run func(ctx context.Context) (int64, error)) *PairDatabase_CountPairs_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllPairs provides a mock function for the type PairDatabase
func (_mock *PairDatabase) GetAllPairs(ctx context.Context) ([]postgres.Pair, error) {
	ret := _mock.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllPairs")
	}

	var r0 []postgres.Pair
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context) ([]postgres.Pair, error)); ok {
		return returnFunc(ctx)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context) []postgres.Pair); ok {
		r0 = returnFunc(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]postgres.Pair)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = returnFunc(ctx)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// PairDatabase_GetAllPairs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllPairs'
type PairDatabase_GetAllPairs_Call struct {
	*mock.Call
}

// GetAllPairs is a helper method to define mock.On call
//   - ctx context.Context
func (_e *PairDatabase_Expecter) GetAllPairs(ctx interface{}) *PairDatabase_GetAllPairs_Call {
	return &PairDatabase_GetAllPairs_Call{Call: _e.mock.On("GetAllPairs", ctx)}
}

func (_c *PairDatabase_GetAllPairs_Call) Run(run func(ctx context.Context)) *PairDatabase_GetAllPairs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		run(
			arg0,
		)
	})
	return _c
}

func (_c *PairDatabase_GetAllPairs_Call) Return(pairs []postgres.Pair, err error) *PairDatabase_GetAllPairs_Call {
	_c.Call.Return(pairs, err)
	return _c
}

func (_c *PairDatabase_GetAllPairs_Call) RunAndReturn(run func(ctx context.Context) ([]postgres.Pair, error)) *PairDatabase_GetAllPairs_Call {
	_c.Call.Return(run)
	return _c
}

// GetPair provides a mock function for the type PairDatabase
func (_mock *PairDatabase) GetPair(ctx context.Context, symbol string) (postgres.Pair, error) {
	ret := _mock.Called(ctx, symbol)

	if len(ret) == 0 {
		panic("no return value specified for GetPair")
	}

	var r0 postgres.Pair
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (postgres.Pair, error)); ok {
		return returnFunc(ctx, symbol)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) postgres.Pair); ok {
		r0 = returnFunc(ctx, symbol)
	} else {
		r0 = ret.Get(0).(postgres.Pair)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, symbol)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// PairDatabase_GetPair_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPair'
type PairDatabase_GetPair_Call struct {
	*mock.Call
}

// GetPair is a helper method to define mock.On call
//   - ctx context.Context
//   - symbol string
func (_e *PairDatabase_Expecter) GetPair(ctx interface{}, symbol interface{}) *PairDatabase_GetPair_Call {
	return &PairDatabase_GetPair_Call{Call: _e.mock.On("GetPair", ctx, symbol)}
}

func (_c *PairDatabase_GetPair_Call) Run(run func(ctx context.Context, symbol string)) *PairDatabase_GetPair_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 string
		if args[1] != nil {
			arg1 = args[1].(string)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *PairDatabase_GetPair_Call) Return(pair postgres.Pair, err error) *PairDatabase_GetPair_Call {
	_c.Call.Return(pair, err)
	return _c
}

func (_c *PairDatabase_GetPair_Call) RunAndReturn(run func(ctx context.Context, symbol string) (postgres.Pair, error)) *PairDatabase_GetPair_Call {
	_c.Call.Return(run)
	return _c
}

// HasPair provides a mock function for the type PairDatabase
func (_mock *PairDatabase) HasPair(ctx context.Context, symbol string) (bool, error) {
	ret := _mock.Called(ctx, symbol)

	if len(ret) == 0 {
		panic("no return value specified for HasPair")
	}

	var r0 bool
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return returnFunc(ctx, symbol)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = returnFunc(ctx, symbol)
	} else {
		r0 = ret.Get(0).(bool)
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = returnFunc(ctx, symbol)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// PairDatabase_HasPair_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasPair'
type PairDatabase_HasPair_Call struct {
	*mock.Call
}

// HasPair is a helper method to define mock.On call
//   - ctx context.Context
//   - symbol string
func (_e *PairDatabase_Expecter) HasPair(ctx interface{}, symbol interface{}) *PairDatabase_HasPair_Call {
	return &PairDatabase_HasPair_Call{Call: _e.mock.On("HasPair", ctx, symbol)}
}

func (_c *PairDatabase_HasPair_Call) Run(run func(ctx context.Context, symbol string)) *PairDatabase_HasPair_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 string
		if args[1] != nil {
			arg1 = args[1].(string)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *PairDatabase_HasPair_Call) Return(b bool, err error) *PairDatabase_HasPair_Call {
	_c.Call.Return(b, err)
	return _c
}

func (_c *PairDatabase_HasPair_Call) RunAndReturn(run func(ctx context.Context, symbol string) (bool, error)) *PairDatabase_HasPair_Call {
	_c.Call.Return(run)
	return _c
}
