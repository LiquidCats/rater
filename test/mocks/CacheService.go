// Code generated by mockery; DO NOT EDIT.
// github.com/vektra/mockery
// template: testify

package mocks

import (
	"context"

	"github.com/LiquidCats/rater/internal/app/domain/entity"
	mock "github.com/stretchr/testify/mock"
)

// NewCacheService creates a new instance of CacheService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCacheService(t interface {
	mock.TestingT
	Cleanup(func())
}) *CacheService {
	mock := &CacheService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

// CacheService is an autogenerated mock type for the CacheService type
type CacheService struct {
	mock.Mock
}

type CacheService_Expecter struct {
	mock *mock.Mock
}

func (_m *CacheService) EXPECT() *CacheService_Expecter {
	return &CacheService_Expecter{mock: &_m.Mock}
}

// GetRate provides a mock function for the type CacheService
func (_mock *CacheService) GetRate(ctx context.Context, pair entity.Pair) (*entity.Rate, error) {
	ret := _mock.Called(ctx, pair)

	if len(ret) == 0 {
		panic("no return value specified for GetRate")
	}

	var r0 *entity.Rate
	var r1 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, entity.Pair) (*entity.Rate, error)); ok {
		return returnFunc(ctx, pair)
	}
	if returnFunc, ok := ret.Get(0).(func(context.Context, entity.Pair) *entity.Rate); ok {
		r0 = returnFunc(ctx, pair)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.Rate)
		}
	}
	if returnFunc, ok := ret.Get(1).(func(context.Context, entity.Pair) error); ok {
		r1 = returnFunc(ctx, pair)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// CacheService_GetRate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRate'
type CacheService_GetRate_Call struct {
	*mock.Call
}

// GetRate is a helper method to define mock.On call
//   - ctx context.Context
//   - pair entity.Pair
func (_e *CacheService_Expecter) GetRate(ctx interface{}, pair interface{}) *CacheService_GetRate_Call {
	return &CacheService_GetRate_Call{Call: _e.mock.On("GetRate", ctx, pair)}
}

func (_c *CacheService_GetRate_Call) Run(run func(ctx context.Context, pair entity.Pair)) *CacheService_GetRate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 entity.Pair
		if args[1] != nil {
			arg1 = args[1].(entity.Pair)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *CacheService_GetRate_Call) Return(rate *entity.Rate, err error) *CacheService_GetRate_Call {
	_c.Call.Return(rate, err)
	return _c
}

func (_c *CacheService_GetRate_Call) RunAndReturn(run func(ctx context.Context, pair entity.Pair) (*entity.Rate, error)) *CacheService_GetRate_Call {
	_c.Call.Return(run)
	return _c
}

// PutRate provides a mock function for the type CacheService
func (_mock *CacheService) PutRate(ctx context.Context, rate entity.Rate) error {
	ret := _mock.Called(ctx, rate)

	if len(ret) == 0 {
		panic("no return value specified for PutRate")
	}

	var r0 error
	if returnFunc, ok := ret.Get(0).(func(context.Context, entity.Rate) error); ok {
		r0 = returnFunc(ctx, rate)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// CacheService_PutRate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PutRate'
type CacheService_PutRate_Call struct {
	*mock.Call
}

// PutRate is a helper method to define mock.On call
//   - ctx context.Context
//   - rate entity.Rate
func (_e *CacheService_Expecter) PutRate(ctx interface{}, rate interface{}) *CacheService_PutRate_Call {
	return &CacheService_PutRate_Call{Call: _e.mock.On("PutRate", ctx, rate)}
}

func (_c *CacheService_PutRate_Call) Run(run func(ctx context.Context, rate entity.Rate)) *CacheService_PutRate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		var arg0 context.Context
		if args[0] != nil {
			arg0 = args[0].(context.Context)
		}
		var arg1 entity.Rate
		if args[1] != nil {
			arg1 = args[1].(entity.Rate)
		}
		run(
			arg0,
			arg1,
		)
	})
	return _c
}

func (_c *CacheService_PutRate_Call) Return(err error) *CacheService_PutRate_Call {
	_c.Call.Return(err)
	return _c
}

func (_c *CacheService_PutRate_Call) RunAndReturn(run func(ctx context.Context, rate entity.Rate) error) *CacheService_PutRate_Call {
	_c.Call.Return(run)
	return _c
}
