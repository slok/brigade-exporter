// Code generated by mockery v1.0.0. DO NOT EDIT.

package brigade

import brigade "github.com/slok/brigade-exporter/pkg/service/brigade"
import mock "github.com/stretchr/testify/mock"

// Interface is an autogenerated mock type for the Interface type
type Interface struct {
	mock.Mock
}

// GetBuilds provides a mock function with given fields:
func (_m *Interface) GetBuilds() ([]*brigade.Build, error) {
	ret := _m.Called()

	var r0 []*brigade.Build
	if rf, ok := ret.Get(0).(func() []*brigade.Build); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*brigade.Build)
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

// GetJobs provides a mock function with given fields:
func (_m *Interface) GetJobs() ([]*brigade.Job, error) {
	ret := _m.Called()

	var r0 []*brigade.Job
	if rf, ok := ret.Get(0).(func() []*brigade.Job); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*brigade.Job)
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

// GetProjects provides a mock function with given fields:
func (_m *Interface) GetProjects() ([]*brigade.Project, error) {
	ret := _m.Called()

	var r0 []*brigade.Project
	if rf, ok := ret.Get(0).(func() []*brigade.Project); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*brigade.Project)
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
