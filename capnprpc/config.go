// Copyright (c) 2025 Matheus Degiovani. All rights reserved.
// Use of this source code is governed by a source-available
// license that can be found in the LICENSE file.

package capnprpc

import (
	"fmt"
	"sync/atomic"

	"github.com/rs/zerolog"
)

type vatConfig struct {
	logger *zerolog.Logger
	name   string
}

// applyOptions applies the given options to the config.
func (vc *vatConfig) applyOptions(opts ...VatOption) {
	for _, opt := range opts {
		opt(vc)
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
