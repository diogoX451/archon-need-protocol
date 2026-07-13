"""Parse need events with dual-read (envelope / bare / legacy top-level).

Mirrors github.com/diogoX451/archon-need-protocol ParseNeedEvent so other
languages can implement against the same golden fixtures under
testdata/golden/.
"""

from __future__ import annotations

import json
from dataclasses import dataclass, field
from typing import Any


@dataclass
class ExternalNeed:
    type: str = ""
    correlation_id: str = ""
    payload: Any = field(default_factory=dict)
    created_at: str = ""


@dataclass
class NeedEvent:
    schema_version: int = 0
    net_id: str = ""
    workflow_id: str = ""
    need: ExternalNeed = field(default_factory=ExternalNeed)


@dataclass
class WebhookBody:
    payload: Any = None


def _as_need(obj: dict[str, Any]) -> ExternalNeed:
    return ExternalNeed(
        type=str(obj.get("type") or ""),
        correlation_id=str(obj.get("correlation_id") or ""),
        payload=obj.get("payload") if obj.get("payload") is not None else {},
        created_at=str(obj.get("created_at") or ""),
    )


def parse_need_event(data: bytes | str | dict[str, Any]) -> NeedEvent:
    """Dual-read ParseNeedEvent compatible with Go goldens."""
    if isinstance(data, (bytes, bytearray)):
        obj = json.loads(data.decode("utf-8"))
    elif isinstance(data, str):
        obj = json.loads(data)
    else:
        obj = data

    # Outer envelope: {schema_version, payload: {...NeedEvent...}}
    if isinstance(obj, dict) and "payload" in obj and "need" not in obj and "correlation_id" not in obj:
        inner = obj["payload"]
        if isinstance(inner, str):
            try:
                inner = json.loads(inner)
            except json.JSONDecodeError:
                pass
        if isinstance(inner, dict) and ("need" in inner or "correlation_id" in inner or "type" in inner):
            obj = inner

    event = NeedEvent(
        schema_version=int(obj.get("schema_version") or 0),
        net_id=str(obj.get("net_id") or ""),
        workflow_id=str(obj.get("workflow_id") or ""),
    )
    if "need" in obj and isinstance(obj["need"], dict):
        event.need = _as_need(obj["need"])
    elif obj.get("correlation_id"):
        # legacy top-level ExternalNeed
        event.need = _as_need(obj)
    return event
