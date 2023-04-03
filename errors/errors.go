// Package errors is used to handle errors
// with ease in golang
package errors

import "net/http"

// Status is the enum to have fixed
// amount of errors that will be forwarded
// to the client
type Status string

// String converts the Status enum to a more
// usable string
func (s Status) String() string {
	return string(s)
}

const (
	// Okay is to represent the succsess message
	Okay Status = "Okay"
	// BadRequest is to represent the Bad request message
	BadRequest Status = "BadRequest"
	// InternalServerError is to represent Internal
	// Server Errors
	InternalServerError Status = "InternalServerError"
	// UnAuthorized is to represent UnAuthorized requests
	UnAuthorized Status = "UnAuthorized"
	// MethodNotAllowed is to represent methods that
	// are not allowed
	MethodNotAllowed Status = "MethodNotAllowed"
	// ContentToLarge is to represent the contents
	// that are too large
	ContentToLarge Status = "ContentToLarge"
	// UnSupportedMedia is to represent content that is not
	// currently supported
	UnSupportedMedia Status = "UnSupportedMedia"
	// CloudVisionFailed is to represent that image does not
	// compliy with the policies
	CloudVisionFailed Status = "CloudVisionFailed"
	// WorkInProgress is to represent that the current route
	// is a work in progress
	WorkInProgress Status = "WorkInProgress"
)

// Response is a struct to the send the
// response to the client
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetStatusCode is a function to get the status
// code for various status enums
func GetStatusCode(s Status) int {
	code := http.StatusOK

	switch s {
	case Okay:
		code = http.StatusOK
		break
	case BadRequest:
		code = http.StatusBadRequest
		break
	case InternalServerError:
		code = http.StatusInternalServerError
		break
	case UnAuthorized:
		code = http.StatusUnauthorized
		break
	case MethodNotAllowed:
		code = http.StatusMethodNotAllowed
		break
	case ContentToLarge:
		code = 415
		break
	case UnSupportedMedia:
		code = http.StatusUnsupportedMediaType
		break
	case CloudVisionFailed:
		code = 418
		break
	case WorkInProgress:
		code = 202
	default:
		code = http.StatusInternalServerError
	}

	return code
}
