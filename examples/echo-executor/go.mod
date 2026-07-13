module github.com/diogoX451/archon-need-protocol/examples/echo-executor

go 1.22

require (
	github.com/diogoX451/archon-executor-runtime v0.1.0
	github.com/diogoX451/archon-nats-bus v0.1.0
	github.com/diogoX451/archon-need-protocol v0.1.0
)

replace github.com/diogoX451/archon-need-protocol => ../..
replace github.com/diogoX451/archon-executor-runtime => ../../../archon-executor-runtime
replace github.com/diogoX451/archon-nats-bus => ../../../archon-nats-bus
replace github.com/diogoX451/archon-event-envelope => ../../../archon-event-envelope
