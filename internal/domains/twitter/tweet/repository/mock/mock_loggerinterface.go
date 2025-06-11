package mock

type MockLoggerInterface struct {
	ErrorFunc func(msg string, keysAndValues ...interface{})
	InfoFunc  func(msg string, keysAndValues ...interface{})
	DebugFunc func(msg string, keysAndValues ...interface{})
}

func (m *MockLoggerInterface) Error(msg string, keysAndValues ...interface{}) {
	if m.ErrorFunc != nil {
		m.ErrorFunc(msg, keysAndValues...)
	}
}

func (m *MockLoggerInterface) Info(msg string, keysAndValues ...interface{}) {
	if m.InfoFunc != nil {
		m.InfoFunc(msg, keysAndValues...)
	}
}

func (m *MockLoggerInterface) Debug(msg string, keysAndValues ...interface{}) {
	if m.DebugFunc != nil {
		m.DebugFunc(msg, keysAndValues...)
	}
}
