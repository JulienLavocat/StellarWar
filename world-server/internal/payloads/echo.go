package payloads

import "encoding/json"

type EchoPayload struct {
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}
