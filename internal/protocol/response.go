package protocol

type Response struct {
	ProtocolVersion int          `json:"protocol_version"`
	RequestID       string       `json:"request_id"`
	Success         bool         `json:"success"`
	Data            any          `json:"data"`
	Error           *ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
