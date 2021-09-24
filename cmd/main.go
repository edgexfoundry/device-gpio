// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides device service of a GPIO devices.
package main

import (
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"
	
	"github.com/edgexfoundry/device-gpio"
	"github.com/edgexfoundry/device-gpio/internal/driver"
)

const (
	serviceName string = "device-gpio"
)

func main() {
	d := driver.Driver{}
	startup.Bootstrap(serviceName, device.Version, &d)
}
