import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "node:test";
import { parseNeedEvent } from "../src/parse.ts";

const goldenDir = join(
  dirname(fileURLToPath(import.meta.url)),
  "../../../testdata/golden",
);

test("golden need_v1", () => {
  const raw = readFileSync(join(goldenDir, "need_v1.json"), "utf8");
  const ev = parseNeedEvent(raw);
  assert.equal(ev.schema_version, 1);
  assert.equal(ev.need.correlation_id, "corr-001");
  assert.equal(ev.need.type, "http");
  assert.equal(
    (ev.need.payload as { url: string }).url,
    "https://example.com/health",
  );
});

test("golden legacy top-level", () => {
  const raw = readFileSync(join(goldenDir, "need_legacy_toplevel.json"), "utf8");
  const ev = parseNeedEvent(raw);
  assert.equal(ev.need.correlation_id, "corr-legacy");
});

test("golden webhook body shape", () => {
  const raw = JSON.parse(
    readFileSync(join(goldenDir, "webhook_body.json"), "utf8"),
  );
  assert.ok(raw.payload);
});

test("correlation may be empty", () => {
  const ev = parseNeedEvent('{"need":{"type":"http","payload":{}}}');
  assert.equal(ev.need.correlation_id, "");
});
