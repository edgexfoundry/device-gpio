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
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func (s *Driver) exportBySysfs(line uint8) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	return ioutil.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d\n", line)), 0644)
}

func (s *Driver) unexportBySysfs(line uint8) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return ioutil.WriteFile("/sys/class/gpio/unexport", []byte(fmt.Sprintf("%d\n", line)), 0644)
	}
	return nil
}

func (s *Driver) setDirectionBySysfs(line uint8, direction string) error {
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
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", line), []byte(way), 0644)
	} else {
		return errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func (s *Driver) getDirectionBySysfs(line uint8) (string, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		direction, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", line))
		if err != nil {
			return "", err
		} else {
			return strings.Replace(string(direction), "\n", "", -1), err
		}
	} else {
		return "", errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func (s *Driver) setValueBySysfs(line uint8, value bool) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var tmp string
		if value {
			tmp = "1"
		} else {
			tmp = "0"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", line), []byte(tmp), 0644)
	} else {
		return errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func (s *Driver) getValueBySysfs(line uint8) (bool, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", line)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ret, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", line))
		if err != nil {
			return false, err
		}
		value, err := strconv.Atoi(strings.Replace(string(ret), "\n", "", -1))
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
