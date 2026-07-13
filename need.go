// Package needprotocol defines the Archon agent need/response wire contracts.
//
// Flow:
//  1. Orchestrator publishes a NeedEvent to a subject (e.g. need.http).
//  2. An external executor consumes it, does work, and POSTs a webhook.
//  3. Orchestrator resumes the paused agent using correlation_id.
//
// Extracted from the Archon agent platform (https://github.com/diogoX451).
package needprotocol

import (
	"encoding/json"
	"time"

	"github.com/diogoX451/archon-event-envelope"
)

// CurrentSchemaVersion mirrors the envelope package.
const CurrentSchemaVersion = envelope.CurrentSchemaVersion

// ExternalNeed is the unit of work requested by a paused agent.
type ExternalNeed struct {
	Type          string          `json:"type"`
	CorrelationID string          `json:"correlation_id"`
	Payload       json.RawMessage `json:"payload"`
	CreatedAt     time.Time       `json:"created_at"`
}

// NeedEvent is the standard bus payload for need messages.
// SchemaVersion is optional (0 = legacy dual-read).
type NeedEvent struct {
	SchemaVersion int          `json:"schema_version,omitempty"`
	NetID         string       `json:"net_id,omitempty"`
	WorkflowID    string       `json:"workflow_id,omitempty"`
	Need          ExternalNeed `json:"need"`
}

// WebhookBody is the JSON body executors POST back to the orchestrator.
// Path convention: POST {base}/webhooks/needs/{correlation_id}
type WebhookBody struct {
	Payload json.RawMessage `json:"payload"`
}

// ParseNeedEvent dual-reads:
//   - outer envelope.Envelope wrapping a NeedEvent
//   - bare NeedEvent JSON
//   - legacy top-level ExternalNeed fields
func ParseNeedEvent(data []byte) (NeedEvent, error) {
	if env, bare, err := envelope.UnmarshalFlexible(data); err == nil && !bare && len(env.Payload) > 0 {
		data = env.Payload
	}
	var event NeedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return NeedEvent{}, err
	}
	if event.Need.CorrelationID == "" && event.NetID == "" {
		var legacy ExternalNeed
		if err := json.Unmarshal(data, &legacy); err == nil && legacy.CorrelationID != "" {
			event.Need = legacy
		}
	}
	return event, nil
}

// StampSchema sets SchemaVersion when unset (legacy 0).
func StampSchema(v *int) {
	if v != nil && *v == 0 {
		*v = CurrentSchemaVersion
	}
}
