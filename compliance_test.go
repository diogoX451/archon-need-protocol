package needprotocol_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	needprotocol "github.com/diogoX451/archon-need-protocol"
)

func TestGolden_NeedV1(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "golden", "need_v1.json"))
	if err != nil {
		t.Fatal(err)
	}
	ev, err := needprotocol.ParseNeedEvent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if ev.SchemaVersion != 1 {
		t.Fatalf("schema=%d", ev.SchemaVersion)
	}
	if ev.Need.CorrelationID != "corr-001" || ev.Need.Type != "http" {
		t.Fatalf("%+v", ev.Need)
	}
	var payload map[string]any
	if err := json.Unmarshal(ev.Need.Payload, &payload); err != nil {
		t.Fatal(err)
	}
	if payload["url"] != "https://example.com/health" {
		t.Fatalf("%v", payload)
	}
}

func TestGolden_LegacyTopLevel(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "golden", "need_legacy_toplevel.json"))
	if err != nil {
		t.Fatal(err)
	}
	ev, err := needprotocol.ParseNeedEvent(raw)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Need.CorrelationID != "corr-legacy" {
		t.Fatalf("%+v", ev)
	}
}

func TestGolden_WebhookBodyShape(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("testdata", "golden", "webhook_body.json"))
	if err != nil {
		t.Fatal(err)
	}
	var body needprotocol.WebhookBody
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatal(err)
	}
	if len(body.Payload) == 0 {
		t.Fatal("empty payload")
	}
}

// Compliance rules that any executor implementation should satisfy.
func TestCompliance_CorrelationIDRequiredForResume(t *testing.T) {
	ev, err := needprotocol.ParseNeedEvent([]byte(`{"need":{"type":"http","payload":{}}}`))
	if err != nil {
		t.Fatal(err)
	}
	if ev.Need.CorrelationID != "" {
		t.Fatal("expected empty correlation for incomplete need")
	}
	// Executors MUST refuse work without correlation_id — documented contract.
}
