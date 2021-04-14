// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides an GPIO implementation of
// ProtocolDriver interface.
//
package driver

import (
	"github.com/warthog618/gpiod"
)

func (s *Driver) getValueByChardev(line *gpiod.Line) (bool, error) {
	val, err := line.Value()
	if err != nil {
		return false, err
	}
	ctx := true
	if val == 0 {
		ctx = false
	}
	return ctx, nil
}

func (s *Driver) setValueByChardev(line *gpiod.Line, ctx int) error {
	if err := line.SetValue(ctx); err != nil {
		return err
	}
	return nil
}
