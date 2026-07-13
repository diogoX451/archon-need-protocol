package needprotocol

import (
	"encoding/json"
	"testing"
	"time"
)

func TestParseNeedEventBare(t *testing.T) {
	t.Parallel()
	raw, _ := json.Marshal(NeedEvent{
		SchemaVersion: 1,
		NetID:         "net-1",
		Need: ExternalNeed{
			Type:          "http",
			CorrelationID: "c-1",
			Payload:       json.RawMessage(`{"url":"https://example.com"}`),
			CreatedAt:     time.Now().UTC(),
		},
	})
	ev, err := ParseNeedEvent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Need.CorrelationID != "c-1" || ev.Need.Type != "http" {
		t.Fatalf("%+v", ev)
	}
}

func TestParseNeedEventLegacyTopLevel(t *testing.T) {
	t.Parallel()
	raw := []byte(`{"type":"http","correlation_id":"c-2","payload":{"a":1}}`)
	ev, err := ParseNeedEvent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Need.CorrelationID != "c-2" {
		t.Fatalf("%+v", ev)
	}
}
