# archon-need-protocol

[![Go Reference](https://pkg.go.dev/badge/github.com/diogoX451/archon-need-protocol.svg)](https://pkg.go.dev/github.com/diogoX451/archon-need-protocol)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

Wire contracts for **agent needs** (pause → external work → webhook resume).

Part of the **Archon open-source toolkit** by [@diogoX451](https://github.com/diogoX451).

## Spec

See [SPEC.md](SPEC.md) for the protocol definition (language-agnostic).

## Install

```bash
go get github.com/diogoX451/archon-need-protocol@latest
```

## Go types

```go
import needprotocol "github.com/diogoX451/archon-need-protocol"

ev, err := needprotocol.ParseNeedEvent(msg.Data())
// POST webhook: needprotocol.WebhookBody{Payload: resultJSON}
```

## Contributing

New language SDKs, golden JSON fixtures, and executor examples are welcome.
See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Apache-2.0
