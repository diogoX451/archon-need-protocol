# archon-need-protocol (TypeScript port)

Minimal dual-read parser for the [Archon need protocol](https://github.com/diogoX451/archon-need-protocol)
goldens under `testdata/golden/`.

## Run golden tests

```bash
cd ports/typescript
npm test
# or without install (Node 20+):
node --experimental-strip-types --test test/golden.test.ts
```

## Usage

```ts
import { parseNeedEvent } from "./src/parse.ts";

const ev = parseNeedEvent(await Deno.readTextFile("../../testdata/golden/need_v1.json"));
console.log(ev.need.correlation_id); // corr-001
```

The Go package remains the source of truth; this port exists so other
languages can implement against the same fixtures.
