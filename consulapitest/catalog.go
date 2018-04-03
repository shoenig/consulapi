// Code autogenerated by mockery v2.0.0
//
// Do not manually edit the content of this file.

// Package consulapitest contains autogenerated mocks.
package consulapitest

import "github.com/shoenig/consulapi"
import "github.com/stretchr/testify/mock"

// Catalog is an autogenerated mock type for the Catalog type
type Catalog struct {
	mock.Mock
}

// Datacenters provides a mock function with given fields:
func (mockerySelf *Catalog) Datacenters() ([]string, error) {
	ret := mockerySelf.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Node provides a mock function with given fields: dc, name
func (mockerySelf *Catalog) Node(dc string, name string) (consulapi.NodeInfo, error) {
	ret := mockerySelf.Called(dc, name)

	var r0 consulapi.NodeInfo
	if rf, ok := ret.Get(0).(func(string, string) consulapi.NodeInfo); ok {
		r0 = rf(dc, name)
	} else {
		r0 = ret.Get(0).(consulapi.NodeInfo)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(dc, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Nodes provides a mock function with given fields: dc
func (mockerySelf *Catalog) Nodes(dc string) ([]consulapi.Node, error) {
	ret := mockerySelf.Called(dc)

	var r0 []consulapi.Node
	if rf, ok := ret.Get(0).(func(string) []consulapi.Node); ok {
		r0 = rf(dc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]consulapi.Node)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Service provides a mock function with given fields: dc, service, tags
func (mockerySelf *Catalog) Service(dc string, service string, tags ...string) ([]consulapi.Service, error) {
	mockeryVariadicArg := make([]interface{}, len(tags))
	for mockeryI := range tags {
		mockeryVariadicArg[mockeryI] = tags[mockeryI]
	}
	var mockeryCalledArg []interface{}
	mockeryCalledArg = append(mockeryCalledArg, dc, service)
	mockeryCalledArg = append(mockeryCalledArg, mockeryVariadicArg...)
	ret := mockerySelf.Called(mockeryCalledArg...)

	var r0 []consulapi.Service
	if rf, ok := ret.Get(0).(func(string, string, ...string) []consulapi.Service); ok {
		r0 = rf(dc, service, tags...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]consulapi.Service)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, ...string) error); ok {
		r1 = rf(dc, service, tags...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Services provides a mock function with given fields: dc
func (mockerySelf *Catalog) Services(dc string) (map[string][]string, error) {
	ret := mockerySelf.Called(dc)

	var r0 map[string][]string
	if rf, ok := ret.Get(0).(func(string) map[string][]string); ok {
		r0 = rf(dc)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string][]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(dc)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
