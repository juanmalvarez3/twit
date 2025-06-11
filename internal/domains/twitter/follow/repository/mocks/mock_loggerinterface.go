package mocks

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLoggerInterface struct {
	mock.Mock
}

func (m *MockLoggerInterface) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLoggerInterface) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLoggerInterface) Warn(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLoggerInterface) Debug(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}
