# archon-need-protocol (Python port)

Minimal dual-read parser for the [Archon need protocol](https://github.com/diogoX451/archon-need-protocol)
goldens under `testdata/golden/`.

## Run golden tests

```bash
cd ports/python
PYTHONPATH=. python -m unittest tests/test_golden.py -v
```

## Usage

```python
from needprotocol import parse_need_event

ev = parse_need_event(open("../../testdata/golden/need_v1.json", "rb").read())
assert ev.need.correlation_id == "corr-001"
```

The Go package remains the source of truth; this port exists so other
languages can implement against the same fixtures.
