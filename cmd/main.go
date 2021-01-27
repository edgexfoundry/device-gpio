// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2020 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides device service of a GPIO devices.
package main

import (
	"github.com/edgexfoundry/device-gpio-go"
	"github.com/edgexfoundry/device-gpio-go/internal/driver"
	"github.com/edgexfoundry/device-sdk-go/pkg/startup"
)

const (
	serviceName string = "device-gpio-go"
)

func main() {
	d := driver.Driver{}
	startup.Bootstrap(serviceName, device.Version, &d)
}
