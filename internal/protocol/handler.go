package protocol

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	"karisdb/internal/app"
)

func Handle(ctx context.Context, input []byte, router *app.Router) Response {
	response := Response{ProtocolVersion: Version}
	if len(input) > MaxRequestBytes {
		return failure(response, "REQUEST_TOO_LARGE", "request exceeds the maximum size")
	}

	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.DisallowUnknownFields()
	var request Request
	if err := decoder.Decode(&request); err != nil {
		return failure(response, "INVALID_REQUEST", "request must be a valid JSON object")
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return failure(response, "INVALID_REQUEST", "input must contain exactly one JSON request")
	}

	response.RequestID = request.RequestID
	if request.ProtocolVersion != Version {
		return failure(response, "UNSUPPORTED_PROTOCOL_VERSION", "protocol version is not supported")
	}
	if strings.TrimSpace(request.RequestID) == "" || len(request.RequestID) > 128 {
		return failure(response, "INVALID_REQUEST_ID", "request_id must contain 1 to 128 characters")
	}
	if request.Command == "" {
		return failure(response, "INVALID_REQUEST", "command is required")
	}

	data, appErr := router.Execute(ctx, request.Command, request.Arguments)
	if appErr != nil {
		return failure(response, appErr.Code, appErr.Message)
	}
	response.Success = true
	response.Data = data
	return response
}

func WriteResponse(writer io.Writer, response Response) error {
	return json.NewEncoder(writer).Encode(response)
}

func failure(response Response, code, message string) Response {
	response.Success = false
	response.Data = nil
	response.Error = &ErrorDetail{Code: code, Message: message}
	return response
}
