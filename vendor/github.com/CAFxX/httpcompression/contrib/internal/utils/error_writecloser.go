package utils

type ErrorWriteCloser struct {
	Err error
}

func (e ErrorWriteCloser) Write(_ []byte) (int, error) {
	return 0, e.Err
}

func (e ErrorWriteCloser) Close() error {
	return e.Err
}
