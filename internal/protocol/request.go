package protocol

import "encoding/json"

const Version = 1
const MaxRequestBytes = 1 << 20

type Request struct {
	ProtocolVersion int                        `json:"protocol_version"`
	RequestID       string                     `json:"request_id"`
	Command         string                     `json:"command"`
	Arguments       map[string]json.RawMessage `json:"arguments"`
}
