"""Golden compliance tests — share fixtures with the Go package."""

from __future__ import annotations

import json
import unittest
from pathlib import Path

from needprotocol.parse import WebhookBody, parse_need_event

# Repo layout: ports/python/tests → ../../../testdata/golden
GOLDEN = Path(__file__).resolve().parents[3] / "testdata" / "golden"


class TestGolden(unittest.TestCase):
    def test_need_v1(self) -> None:
        raw = (GOLDEN / "need_v1.json").read_bytes()
        ev = parse_need_event(raw)
        self.assertEqual(ev.schema_version, 1)
        self.assertEqual(ev.need.correlation_id, "corr-001")
        self.assertEqual(ev.need.type, "http")
        self.assertEqual(ev.need.payload["url"], "https://example.com/health")

    def test_legacy_toplevel(self) -> None:
        raw = (GOLDEN / "need_legacy_toplevel.json").read_bytes()
        ev = parse_need_event(raw)
        self.assertEqual(ev.need.correlation_id, "corr-legacy")

    def test_webhook_body_shape(self) -> None:
        raw = (GOLDEN / "webhook_body.json").read_bytes()
        obj = json.loads(raw)
        body = WebhookBody(payload=obj.get("payload"))
        self.assertIsNotNone(body.payload)
        self.assertTrue(body.payload)

    def test_correlation_may_be_empty(self) -> None:
        ev = parse_need_event(b'{"need":{"type":"http","payload":{}}}')
        self.assertEqual(ev.need.correlation_id, "")


if __name__ == "__main__":
    unittest.main()
