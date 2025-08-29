// This is a separate module in order to avoid having the grpc dependency on the
// main go module.
module matheusd.com/mdcapnp/internal/experiments/echogrpcserver

go 1.24.3

require (
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.8
	matheusd.com/depvendoredtestify v1.10.0-alpha
)

require (
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)
