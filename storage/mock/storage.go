// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package storagemock is a generated GoMock package.
package storagemock

import (
	context "context"
	reflect "reflect"

	model "github.com/ale0x78ey/yandex-practicum-go-developer-devops/model"
	gomock "github.com/golang/mock/gomock"
)

// MockMetricStorage is a mock of MetricStorage interface.
type MockMetricStorage struct {
	ctrl     *gomock.Controller
	recorder *MockMetricStorageMockRecorder
}

// MockMetricStorageMockRecorder is the mock recorder for MockMetricStorage.
type MockMetricStorageMockRecorder struct {
	mock *MockMetricStorage
}

// NewMockMetricStorage creates a new mock instance.
func NewMockMetricStorage(ctrl *gomock.Controller) *MockMetricStorage {
	mock := &MockMetricStorage{ctrl: ctrl}
	mock.recorder = &MockMetricStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricStorage) EXPECT() *MockMetricStorageMockRecorder {
	return m.recorder
}

// Flush mocks base method.
func (m *MockMetricStorage) Flush(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Flush", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Flush indicates an expected call of Flush.
func (mr *MockMetricStorageMockRecorder) Flush(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Flush", reflect.TypeOf((*MockMetricStorage)(nil).Flush), ctx)
}

// IncrMetric mocks base method.
func (m *MockMetricStorage) IncrMetric(ctx context.Context, metric model.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrMetric", ctx, metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrMetric indicates an expected call of IncrMetric.
func (mr *MockMetricStorageMockRecorder) IncrMetric(ctx, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrMetric", reflect.TypeOf((*MockMetricStorage)(nil).IncrMetric), ctx, metric)
}

// LoadMetric mocks base method.
func (m *MockMetricStorage) LoadMetric(ctx context.Context, metricType model.MetricType, metricName string) (*model.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadMetric", ctx, metricType, metricName)
	ret0, _ := ret[0].(*model.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadMetric indicates an expected call of LoadMetric.
func (mr *MockMetricStorageMockRecorder) LoadMetric(ctx, metricType, metricName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadMetric", reflect.TypeOf((*MockMetricStorage)(nil).LoadMetric), ctx, metricType, metricName)
}

// LoadMetricList mocks base method.
func (m *MockMetricStorage) LoadMetricList(ctx context.Context) ([]model.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadMetricList", ctx)
	ret0, _ := ret[0].([]model.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadMetricList indicates an expected call of LoadMetricList.
func (mr *MockMetricStorageMockRecorder) LoadMetricList(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadMetricList", reflect.TypeOf((*MockMetricStorage)(nil).LoadMetricList), ctx)
}

// SaveMetric mocks base method.
func (m *MockMetricStorage) SaveMetric(ctx context.Context, metric model.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMetric", ctx, metric)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMetric indicates an expected call of SaveMetric.
func (mr *MockMetricStorageMockRecorder) SaveMetric(ctx, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetric", reflect.TypeOf((*MockMetricStorage)(nil).SaveMetric), ctx, metric)
}
