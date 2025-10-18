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

type vatConfig struct {
	logger           *zerolog.Logger
	name             string
	bootstrapHandler CallHandler
	net              VatNetwork
	failOnConnErr    bool

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
	if vc.bootstrapHandler == nil {
		vc.bootstrapHandler = allUnimplementedCallHandler{}
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

func WithName(name string) VatOption {
	return func(c *vatConfig) {
		c.name = name
	}
}

func WithLogger(l *zerolog.Logger) VatOption {
	return func(c *vatConfig) {
		c.logger = l
	}
}

func WithVatNetwork(net VatNetwork) VatOption {
	return func(c *vatConfig) {
		c.net = net
	}
}

func WithBootstrapHandler(h CallHandler) VatOption {
	return func(c *vatConfig) {
		c.bootstrapHandler = h
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
