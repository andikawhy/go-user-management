package helper

type StandardError struct {
	Error     error `json:"error"`
	ErrorCode uint  `json:"error_code"`
}
