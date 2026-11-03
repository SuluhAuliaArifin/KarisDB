package protocol

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"karisdb/internal/app"
)

func newTestRouter(t *testing.T) *app.Router {
	t.Helper()
	router, err := app.NewRouter()
	if err != nil {
		t.Fatal(err)
	}
	return router
}

func TestHandleVersionRequest(t *testing.T) {
	request := `{"protocol_version":1,"request_id":"test-1","command":"version","arguments":{}}`
	response := Handle(context.Background(), []byte(request), newTestRouter(t))
	if !response.Success {
		t.Fatalf("request failed: %#v", response.Error)
	}
	if response.RequestID != "test-1" {
		t.Fatalf("RequestID = %q", response.RequestID)
	}
	data, ok := response.Data.(map[string]string)
	if !ok || data["version"] != app.Version {
		t.Fatalf("unexpected response data: %#v", response.Data)
	}
}

func TestHandleRejectsUnsupportedProtocol(t *testing.T) {
	request := `{"protocol_version":2,"request_id":"test-2","command":"version","arguments":{}}`
	response := Handle(context.Background(), []byte(request), newTestRouter(t))
	if response.Success || response.Error == nil || response.Error.Code != "UNSUPPORTED_PROTOCOL_VERSION" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestHandleRejectsOversizedRequest(t *testing.T) {
	response := Handle(context.Background(), []byte(strings.Repeat("x", MaxRequestBytes+1)), newTestRouter(t))
	if response.Success || response.Error == nil || response.Error.Code != "REQUEST_TOO_LARGE" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestHandleRejectsUnknownRequestFields(t *testing.T) {
	request := `{"protocol_version":1,"request_id":"test-3","command":"version","arguments":{},"extra":true}`
	response := Handle(context.Background(), []byte(request), newTestRouter(t))
	if response.Success || response.Error == nil || response.Error.Code != "INVALID_REQUEST" {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestWriteResponseIsOneJSONLine(t *testing.T) {
	var output strings.Builder
	response := Response{ProtocolVersion: Version, RequestID: "test-4", Success: true, Data: json.RawMessage(`{"ok":true}`)}
	if err := WriteResponse(&output, response); err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(output.String(), "\n") {
		t.Fatalf("response is not newline-terminated: %q", output.String())
	}
	var decoded Response
	if err := json.Unmarshal([]byte(output.String()), &decoded); err != nil {
		t.Fatalf("response is invalid JSON: %v", err)
	}
}
