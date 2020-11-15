package error

const (
	ClientError = 0
	ServerError = 1
)

type RequestError struct {
	errType int
	err     error
}

func NewError(errType int, err error) RequestError {
	return RequestError{
		errType: errType,
		err:     err,
	}
}

func (r RequestError) Error() string {
	return r.err.Error()
}

func (r RequestError) IsClientError() bool {
	return r.errType == ClientError
}

func (r RequestError) IsServerError() bool {
	return r.errType == ServerError
}
