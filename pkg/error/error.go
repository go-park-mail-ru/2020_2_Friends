package error

import (
	"net/http"
)

const (
	ClientError = 0
	ServerError = 1
)

type RequestError struct {
	errType int
	err     error
}

func NewClientError(err error) RequestError {
	return RequestError{
		errType: ClientError,
		err:     err,
	}
}

func NewServerError(err error) RequestError {
	return RequestError{
		errType: ServerError,
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

func HandleErrorAndWriteResponse(w http.ResponseWriter, err error, clientErrorStatusCode int) {
	re, ok := err.(RequestError)
	if ok {
		if re.IsClientError() {
			w.WriteHeader(clientErrorStatusCode)
			return
		}
		if re.IsServerError() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
}
