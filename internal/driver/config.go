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
	"fmt"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
)

type configuration struct {
	Abi_driver    string
	Chip_selected string
}

const (
	ABI_DRIVER    = "Interface"
	CHIP_SELECTED = "ChipSelected"
)

func loadInterfaceConfig() (*configuration, error) {
	config := new(configuration)
	if val, ok := service.DriverConfigs()[ABI_DRIVER]; ok {
		config.Abi_driver = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", ABI_DRIVER)
	}
	if val, ok := service.DriverConfigs()[CHIP_SELECTED]; ok {
		config.Chip_selected = val
	} else {
		return config, fmt.Errorf("driver config undefined: %s", CHIP_SELECTED)
	}
	return config, nil
}
