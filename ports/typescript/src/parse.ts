/**
 * Dual-read ParseNeedEvent compatible with Go goldens under testdata/golden/.
 * Mirrors github.com/diogoX451/archon-need-protocol.
 */

export interface ExternalNeed {
  type: string;
  correlation_id: string;
  payload: unknown;
  created_at?: string;
}

export interface NeedEvent {
  schema_version: number;
  net_id: string;
  workflow_id: string;
  need: ExternalNeed;
}

export interface WebhookBody {
  payload: unknown;
}

function asNeed(obj: Record<string, unknown>): ExternalNeed {
  return {
    type: String(obj.type ?? ""),
    correlation_id: String(obj.correlation_id ?? ""),
    payload: obj.payload ?? {},
    created_at: obj.created_at ? String(obj.created_at) : undefined,
  };
}

export function parseNeedEvent(data: string | Uint8Array | object): NeedEvent {
  let obj: Record<string, unknown>;
  if (typeof data === "string") {
    obj = JSON.parse(data) as Record<string, unknown>;
  } else if (data instanceof Uint8Array) {
    obj = JSON.parse(new TextDecoder().decode(data)) as Record<string, unknown>;
  } else {
    obj = data as Record<string, unknown>;
  }

  // Outer envelope: { schema_version, payload: NeedEvent }
  if (
    obj &&
    typeof obj === "object" &&
    "payload" in obj &&
    !("need" in obj) &&
    !("correlation_id" in obj)
  ) {
    let inner = obj.payload;
    if (typeof inner === "string") {
      try {
        inner = JSON.parse(inner);
      } catch {
        /* keep */
      }
    }
    if (
      inner &&
      typeof inner === "object" &&
      ("need" in (inner as object) ||
        "correlation_id" in (inner as object) ||
        "type" in (inner as object))
    ) {
      obj = inner as Record<string, unknown>;
    }
  }

  const event: NeedEvent = {
    schema_version: Number(obj.schema_version ?? 0),
    net_id: String(obj.net_id ?? ""),
    workflow_id: String(obj.workflow_id ?? ""),
    need: { type: "", correlation_id: "", payload: {} },
  };

  if (obj.need && typeof obj.need === "object") {
    event.need = asNeed(obj.need as Record<string, unknown>);
  } else if (obj.correlation_id) {
    event.need = asNeed(obj);
  }
  return event;
}
