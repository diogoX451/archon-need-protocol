# archon-need-protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/diogoX451/archon-need-protocol.svg)](https://pkg.go.dev/github.com/diogoX451/archon-need-protocol)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

**Spec + golden suite + Go types** for the agent **need protocol**: how an
orchestrator offloads work to external executors and resumes by `correlation_id`.

This is the piece the community can implement against in any language.

## What you get

| Artifact | Value |
|----------|--------|
| [SPEC.md](SPEC.md) | Language-agnostic protocol |
| `testdata/golden/` | Fixtures for multi-language ports |
| `compliance_test.go` | Go compliance tests |
| `ports/typescript/` | TS dual-read parser + golden tests |
| `ports/python/` | Python dual-read parser + golden tests |
| `examples/echo-executor` | End-to-end reference worker |
| [archon-http-executor](https://github.com/diogoX451/archon-http-executor) | HTTP reference need worker |
| `ParseNeedEvent` | Dual-read v1 + legacy |

## Why not “just use Celery / Temporal?”

Those are full workflow engines. The need protocol is a **thin pause/resume
contract** you can bolt onto an existing runtime (Interaction Nets, actors,
or a simple state machine) without adopting a new platform.

## Install

```bash
go get github.com/diogoX451/archon-need-protocol@latest
```

## Language ports (same goldens)

```bash
# TypeScript (Node 20+)
cd ports/typescript && node --experimental-strip-types --test test/golden.test.ts

# Python 3.10+
cd ports/python && PYTHONPATH=. python -m unittest tests/test_golden.py -v
```

## High-value contributions

1. More golden cases (envelope-wrapped needs, empty payload, unicode)  
2. Reference executors: shell (sandboxed), MCP bridge  
3. AsyncAPI YAML generation from the SPEC  
4. Publish `@diogox451/archon-need-protocol` to npm / PyPI  

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Related

- [archon-executor-runtime](https://github.com/diogoX451/archon-executor-runtime)  
- [archon-nats-bus](https://github.com/diogoX451/archon-nats-bus)  
- [archon-oss](https://github.com/diogoX451/archon-oss)  

## License

Apache-2.0
