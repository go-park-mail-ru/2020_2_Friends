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

func HandleErrorAndWriteResponse(w http.ResponseWriter, err error) {
	re, ok := err.(RequestError)
	if ok {
		if re.IsClientError() {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if re.IsServerError() {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
}
