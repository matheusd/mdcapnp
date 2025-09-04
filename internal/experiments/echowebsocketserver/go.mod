// This is a separate module in order to avoid having websocket dependencies on
// the main go module.
module matheusd.com/mdcapnp/internal/experiments/echowebsocketserver

go 1.24.3

require (
	github.com/gorilla/websocket v1.5.3
	matheusd.com/testctx v0.2.0
)
