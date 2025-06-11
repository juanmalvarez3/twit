package logger

func ProvideInfo() (*Logger, error) {
	return New("info", "stage")
}

func ProvideWarn() (*Logger, error) {
	return New("warn", "stage")
}

func ProvideError() (*Logger, error) {
	return New("error", "stage")
}
