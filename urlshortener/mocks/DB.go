// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// DB is an autogenerated mock type for the DB type
type DB struct {
	mock.Mock
}

type DB_Expecter struct {
	mock *mock.Mock
}

func (_m *DB) EXPECT() *DB_Expecter {
	return &DB_Expecter{mock: &_m.Mock}
}

// GetURL provides a mock function with given fields: key
func (_m *DB) GetURL(key string) (string, error) {
	ret := _m.Called(key)

	if len(ret) == 0 {
		panic("no return value specified for GetURL")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(key)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DB_GetURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetURL'
type DB_GetURL_Call struct {
	*mock.Call
}

// GetURL is a helper method to define mock.On call
//   - key string
func (_e *DB_Expecter) GetURL(key interface{}) *DB_GetURL_Call {
	return &DB_GetURL_Call{Call: _e.mock.On("GetURL", key)}
}

func (_c *DB_GetURL_Call) Run(run func(key string)) *DB_GetURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *DB_GetURL_Call) Return(_a0 string, _a1 error) *DB_GetURL_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *DB_GetURL_Call) RunAndReturn(run func(string) (string, error)) *DB_GetURL_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function with no fields
func (_m *DB) Init() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Init")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DB_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type DB_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
func (_e *DB_Expecter) Init() *DB_Init_Call {
	return &DB_Init_Call{Call: _e.mock.On("Init")}
}

func (_c *DB_Init_Call) Run(run func()) *DB_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *DB_Init_Call) Return(_a0 error) *DB_Init_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DB_Init_Call) RunAndReturn(run func() error) *DB_Init_Call {
	_c.Call.Return(run)
	return _c
}

// StoreURL provides a mock function with given fields: url, key
func (_m *DB) StoreURL(url string, key string) error {
	ret := _m.Called(url, key)

	if len(ret) == 0 {
		panic("no return value specified for StoreURL")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(url, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DB_StoreURL_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'StoreURL'
type DB_StoreURL_Call struct {
	*mock.Call
}

// StoreURL is a helper method to define mock.On call
//   - url string
//   - key string
func (_e *DB_Expecter) StoreURL(url interface{}, key interface{}) *DB_StoreURL_Call {
	return &DB_StoreURL_Call{Call: _e.mock.On("StoreURL", url, key)}
}

func (_c *DB_StoreURL_Call) Run(run func(url string, key string)) *DB_StoreURL_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *DB_StoreURL_Call) Return(_a0 error) *DB_StoreURL_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DB_StoreURL_Call) RunAndReturn(run func(string, string) error) *DB_StoreURL_Call {
	_c.Call.Return(run)
	return _c
}

// NewDB creates a new instance of DB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDB(t interface {
	mock.TestingT
	Cleanup(func())
}) *DB {
	mock := &DB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
