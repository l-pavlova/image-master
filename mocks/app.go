// Code generated by MockGen. DO NOT EDIT.
// Source: D:/go/src/GoCourse/image-master/image-master/app/app.go

// Package mock_app is a generated GoMock package.
package mock_app

import (
	image "image"
	color "image/color"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	mongo "github.com/l-pavlova/image-master/mongo"
	tensorflowAPI "github.com/l-pavlova/image-master/tensorflowAPI"
)

// MockChangeable is a mock of Changeable interface.
type MockChangeable struct {
	ctrl     *gomock.Controller
	recorder *MockChangeableMockRecorder
}

// MockChangeableMockRecorder is the mock recorder for MockChangeable.
type MockChangeableMockRecorder struct {
	mock *MockChangeable
}

// NewMockChangeable creates a new mock instance.
func NewMockChangeable(ctrl *gomock.Controller) *MockChangeable {
	mock := &MockChangeable{ctrl: ctrl}
	mock.recorder = &MockChangeableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChangeable) EXPECT() *MockChangeableMockRecorder {
	return m.recorder
}

// Set mocks base method.
func (m *MockChangeable) Set(x, y int, c color.Color) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Set", x, y, c)
}

// Set indicates an expected call of Set.
func (mr *MockChangeableMockRecorder) Set(x, y, c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockChangeable)(nil).Set), x, y, c)
}

// MockMongoClient is a mock of MongoClient interface.
type MockMongoClient struct {
	ctrl     *gomock.Controller
	recorder *MockMongoClientMockRecorder
}

// MockMongoClientMockRecorder is the mock recorder for MockMongoClient.
type MockMongoClientMockRecorder struct {
	mock *MockMongoClient
}

// NewMockMongoClient creates a new mock instance.
func NewMockMongoClient(ctrl *gomock.Controller) *MockMongoClient {
	mock := &MockMongoClient{ctrl: ctrl}
	mock.recorder = &MockMongoClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMongoClient) EXPECT() *MockMongoClientMockRecorder {
	return m.recorder
}

// AddImageClassification mocks base method.
func (m *MockMongoClient) AddImageClassification(imagePath string, probabilities []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddImageClassification", imagePath, probabilities)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddImageClassification indicates an expected call of AddImageClassification.
func (mr *MockMongoClientMockRecorder) AddImageClassification(imagePath, probabilities interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddImageClassification", reflect.TypeOf((*MockMongoClient)(nil).AddImageClassification), imagePath, probabilities)
}

// GetAllImageClassifications mocks base method.
func (m *MockMongoClient) GetAllImageClassifications() ([]mongo.ImageClassification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllImageClassifications")
	ret0, _ := ret[0].([]mongo.ImageClassification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllImageClassifications indicates an expected call of GetAllImageClassifications.
func (mr *MockMongoClientMockRecorder) GetAllImageClassifications() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllImageClassifications", reflect.TypeOf((*MockMongoClient)(nil).GetAllImageClassifications))
}

// GetImageClassification mocks base method.
func (m *MockMongoClient) GetImageClassification(imagePath string) (mongo.ImageClassification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetImageClassification", imagePath)
	ret0, _ := ret[0].(mongo.ImageClassification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImageClassification indicates an expected call of GetImageClassification.
func (mr *MockMongoClientMockRecorder) GetImageClassification(imagePath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImageClassification", reflect.TypeOf((*MockMongoClient)(nil).GetImageClassification), imagePath)
}

// MockTensorFlowClient is a mock of TensorFlowClient interface.
type MockTensorFlowClient struct {
	ctrl     *gomock.Controller
	recorder *MockTensorFlowClientMockRecorder
}

// MockTensorFlowClientMockRecorder is the mock recorder for MockTensorFlowClient.
type MockTensorFlowClientMockRecorder struct {
	mock *MockTensorFlowClient
}

// NewMockTensorFlowClient creates a new mock instance.
func NewMockTensorFlowClient(ctrl *gomock.Controller) *MockTensorFlowClient {
	mock := &MockTensorFlowClient{ctrl: ctrl}
	mock.recorder = &MockTensorFlowClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTensorFlowClient) EXPECT() *MockTensorFlowClientMockRecorder {
	return m.recorder
}

// ClassifyImage mocks base method.
func (m *MockTensorFlowClient) ClassifyImage(image image.Image) ([]tensorflowAPI.Label, []float32, [][]float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClassifyImage", image)
	ret0, _ := ret[0].([]tensorflowAPI.Label)
	ret1, _ := ret[1].([]float32)
	ret2, _ := ret[2].([][]float32)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ClassifyImage indicates an expected call of ClassifyImage.
func (mr *MockTensorFlowClientMockRecorder) ClassifyImage(image interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClassifyImage", reflect.TypeOf((*MockTensorFlowClient)(nil).ClassifyImage), image)
}
