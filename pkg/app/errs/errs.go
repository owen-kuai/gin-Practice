package errs

import (
	"fmt"
	"net/http"
)

type PraticeException struct {
	Code      int    `json:"status_code"`
	ErrorCode string `json:"error_code"`
	Msg       string `json:"msg"`
	Details   string `json:"detail"`
	Request   string `json:"request"`
}

func New(errCode, details string) *PraticeException {
	e := &PraticeException{
		ErrorCode: errCode,
		Msg:       codeMap[errCode].message,
		Details:   details,
	}
	e.Code = e.StatusCode()
	return e
}

func (ne *PraticeException) Error() string {
	return fmt.Sprintf("status code: %d, Code: %s, Message: %s, Detail: %s", ne.Code, ne.ErrorCode, ne.Msg, ne.Details)
}

func (ne *PraticeException) StatusCode() int {
	return codeMap[ne.ErrorCode].status_code
}

func Is(err error, code string) bool {
	if err == nil {
		return false
	}

	PraticeErr, ok := err.(*PraticeException)
	if !ok {
		return false
	}

	return PraticeErr.ErrorCode == code
}

const (
	INTERNAL_ERROR = "IE001"
	UNKNOWN_ERROR  = "UK001"

	BADREQUEST_ERROR = "BR001"
)

var codeMap = map[string]struct {
	message     string
	status_code int
}{
	INTERNAL_ERROR:   {message: "internal erorr", status_code: http.StatusInternalServerError},
	UNKNOWN_ERROR:    {message: "unknown erorr", status_code: http.StatusBadRequest},
	BADREQUEST_ERROR: {message: "api error: bad request", status_code: http.StatusBadRequest},
}

func UnknownError(message string) *PraticeException {
	return New(UNKNOWN_ERROR, message)
}

func InternalServerError(message string) *PraticeException {
	return New(INTERNAL_ERROR, message)
}

func ApiBadRequestError(message string) *PraticeException {
	return New(BADREQUEST_ERROR, message)
}
