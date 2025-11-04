// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Jiangxing Intelligence Ltd
//
// SPDX-License-Identifier: Apache-2.0

// Package driver this package provides an GPIO implementation of
// ProtocolDriver interface.
package driver

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (s *Driver) exportBySysfs(line uint16) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	return os.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d\n", line)), 0644) //nolint:gosec
}

func (s *Driver) unexportBySysfs(line uint16) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return os.WriteFile("/sys/class/gpio/unexport", []byte(fmt.Sprintf("%d\n", line)), 0644) //nolint:gosec
	}
	return nil
}

func (s *Driver) setDirectionBySysfs(line uint16, direction string) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var way string
		switch direction {
		case "in":
			way = "in"
		case "out":
			way = "out"
		default:
			return errors.New("invalid direction")
		}
		return os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", line), []byte(way), 0644) //nolint:gosec
	} else {
		return errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func (s *Driver) setValueBySysfs(line uint16, value bool) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var tmp string
		if value {
			tmp = "1"
		} else {
			tmp = "0"
		}
		return os.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", line), []byte(tmp), 0644) //nolint:gosec
	} else {
		return errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func (s *Driver) getValueBySysfs(line uint16) (bool, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ret, err := os.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", line))
		if err != nil {
			return false, err
		}
		value, err := strconv.Atoi(strings.ReplaceAll(string(ret), "\n", ""))
		if err != nil {
			return false, err
		}
		switch value {
		case 1:
			return true, nil
		case 0:
			return false, nil
		default:
			return false, errors.New("invalid value")
		}
	} else {
		return false, errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}
