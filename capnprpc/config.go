// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

// vatConfig holds the full config of a [Vat].
type vatConfig struct {
	logger        *zerolog.Logger
	name          string
	bootServer    CapabilityServer
	net           VatNetwork
	failOnConnErr bool

	// delayResolveIn3PH is how long to delay sending a resolve after
	// sending a Provide.
	delayResolveIn3PH time.Duration
}

// applyOptions applies the given options to the config.
func (vc *vatConfig) applyOptions(opts ...VatOption) {
	for _, opt := range opts {
		opt(vc)
	}

	// If no forms of a bootstrap cap were specified, use a fixed one that
	// returns everything as unimplemented.
	if vc.bootServer == nil {
		vc.bootServer = allUnimplementedCallHandler{}
	}
}

// vatLogger returns the logger to use with the vat.
func (vc *vatConfig) vatLogger() *zerolog.Logger {
	logger := vc.logger
	if logger != nil {
		l := logger.With().Str("vat", vc.name).Logger()
		logger = &l
	} else {
		l := zerolog.Nop()
		logger = &l
	}

	return logger
}

// totalVats tracks the total number of vats.
var totalVats atomic.Uint32

// defaultVatConfig returns the default vat config.
func defaultVatConfig() vatConfig {
	return vatConfig{
		name: fmt.Sprintf("vat%08x", totalVats.Add(1)),
	}
}

type VatOption func(c *vatConfig)

// WithName defines the name to use for a [Vat] instance. This is mostly used
// for debugging/logging purposes.
func WithName(name string) VatOption {
	return func(c *vatConfig) {
		c.name = name
	}
}

// WithLogger defines the logger to use with a [Vat] instance.
func WithLogger(l *zerolog.Logger) VatOption {
	return func(c *vatConfig) {
		c.logger = l
	}
}

// WithVatNetwork defines the [VatNetwork] that a particular [Vat] belongs to.
// This is only necessary for level 3+ vats that intend to use third party
// capability introduction (3PH).
func WithVatNetwork(net VatNetwork) VatOption {
	return func(c *vatConfig) {
		c.net = net
	}
}

// WithBootstrapServer defines the [CapabilityServer] that is to be used as the
// bootstrap capability for remote vats that connect to this vat.
//
// If no bootstrap server is set on a [Vat], then every call made by the remote
// will return an error.
func WithBootstrapServer(cs CapabilityServer) VatOption {
	return func(c *vatConfig) {
		c.bootServer = cs
	}
}

// withFailOnConnErr sets up the vat to fail its own Run() function if any
// connection fails for any reason other than a graceful close.
func withFailOnConnErr(fail bool) VatOption {
	return func(c *vatConfig) {
		c.failOnConnErr = fail
	}
}

// withDelayResolveIn3PH sets up how long to wait after sending a Provide
// message to send the corresponding Resolve in a 3PH scenario.
func withDelayResolveIn3PH(d time.Duration) VatOption {
	return func(c *vatConfig) {
		c.delayResolveIn3PH = d
	}
}
