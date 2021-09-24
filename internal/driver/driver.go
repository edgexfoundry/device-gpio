// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018-2019 IOTech Ltd
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
	"time"

	dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/common"

	"github.com/spf13/cast"
	"github.com/warthog618/gpiod"
)

type Driver struct {
	lc         logger.LoggingClient
	asyncCh    chan<- *dsModels.AsyncValues
	config     *configuration
	openedChip *gpiod.Chip
	openedLine map[string](*gpiod.Line)
}

// Initialize performs protocol-specific initialization for the device
// service.
func (s *Driver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {
	s.lc = lc
	s.asyncCh = asyncCh

	s.openedChip = nil
	s.openedLine = make(map[string](*gpiod.Line))

	config, err := loadInterfaceConfig()
	if err != nil {
		panic(fmt.Errorf("load GPIO configuration failed: %v", err))
	}
	if config.Abi_driver != "sysfs" && config.Abi_driver != "chardev" {
		panic(fmt.Errorf("unsupport GPIO ABI interface: %v", config.Abi_driver))
	}
	s.config = config
	s.lc.Info(fmt.Sprintf("Interface: %v", config.Abi_driver))
	s.lc.Info(fmt.Sprintf("ChipSelected: %v", config.Chip_selected))
	return nil
}

// HandleReadCommands triggers a protocol Read operation for the specified device.
func (s *Driver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {

	s.lc.Info(fmt.Sprintf("protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))
	// fmt.Println(fmt.Sprintf("protocols: %v resource: %v attributes: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes))
	if s.openedChip == nil && s.config.Abi_driver == "chardev" {
		valid_chip, err := cast.ToUint8E(s.config.Chip_selected)
		if err != nil {
			s.lc.Error("invalid chip number, override with gpiochip0, %v", err)
			valid_chip = 0
		}
		chipName := fmt.Sprintf("gpiochip%d", valid_chip)
		s.openedChip, err = gpiod.NewChip(chipName)
		if err != nil {
			s.lc.Error("failed to open %v, %v", chipName, err)
		}
	}
	lineNumStr := fmt.Sprintf("%v", reqs[0].Attributes["line"])

	// now := time.Now().UnixNano()

	for _, req := range reqs {
		value, err := s.getGPIO(lineNumStr)
		if err != nil {
			return nil, err
		}
		cv, _ := dsModels.NewCommandValue(req.DeviceResourceName, common.ValueTypeBool, value)
		res = append(res, cv)
	}
	return res, nil
}

// HandleWriteCommands passes a slice of CommandRequest struct each representing
// a ResourceOperation for a specific device resource.
// Since the commands are actuation commands, params provide parameters for the individual
// command.
func (s *Driver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {
	s.lc.Info(fmt.Sprintf("Driver.HandleWriteCommands: protocols: %v, resource: %v, attribute: %v, parameters: %v", protocols, reqs[0].DeviceResourceName, reqs[0].Attributes, params))
	if s.openedChip == nil && s.config.Abi_driver == "chardev" {
		valid_chip, err := cast.ToUint8E(s.config.Chip_selected)
		if err != nil {
			s.lc.Error("invalid chip number, override with gpiochip0, %v", err)
			valid_chip = 0
		}
		chipName := fmt.Sprintf("gpiochip%d", valid_chip)
		s.openedChip, err = gpiod.NewChip(chipName)
		if err != nil {
			s.lc.Error("failed to open %v, %v", chipName, err)
		}
	}
	lineNumStr := fmt.Sprintf("%v", reqs[0].Attributes["line"])

	for _, param := range params {
		value, err := param.BoolValue()
		if err != nil {
			return err
		}
		if err := s.setGPIO(lineNumStr, value); err != nil {
			return err
		}
	}
	return nil
}

// Stop the protocol-specific DS code to shutdown gracefully, or
// if the force parameter is 'true', immediately. The driver is responsible
// for closing any in-use channels, including the channel used to send async
// readings (if supported).
func (s *Driver) Stop(force bool) error {
	// Then Logging Client might not be initialized
	if s.lc != nil {
		s.lc.Debug(fmt.Sprintf("Driver.Stop called: force=%v", force))
	}
	switch s.config.Abi_driver {
	case "sysfs":
		{
			for line := range s.openedLine {
				valid_line, err := cast.ToUint8E(line)
				if err != nil {
					s.lc.Debug(fmt.Sprintf("Driver.Stop: invalid line %v", line))
					continue
				}
				if err := s.unexportBySysfs(valid_line); err != nil {
					s.lc.Debug(fmt.Sprintf("Driver.Stop: failed to unexport %v", line))
					continue
				}
			}
		}
	case "chardev":
		{
			for line, chip := range s.openedLine {
				if err := chip.Close(); err != nil {
					s.lc.Debug(fmt.Sprintf("Driver.Stop: failed to close line %v, %v", line, err))
					continue
				}
			}
			if err := s.openedChip.Close(); err != nil {
				s.lc.Debug(fmt.Sprintf("Driver.Stop: failed to close chip, %v", err))
			}
		}
	}
	return nil
}

// AddDevice is a callback function that is invoked
// when a new Device associated with this Device Service is added
func (s *Driver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	s.lc.Debug(fmt.Sprintf("a new Device is added: %s", deviceName))
	return nil
}

// UpdateDevice is a callback function that is invoked
// when a Device associated with this Device Service is updated
func (s *Driver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	s.lc.Debug(fmt.Sprintf("Device %s is updated", deviceName))
	return nil
}

// RemoveDevice is a callback function that is invoked
// when a Device associated with this Device Service is removed
func (s *Driver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	s.lc.Debug(fmt.Sprintf("Device %s is removed", deviceName))
	return nil
}

//
func (s *Driver) getGPIO(line string) (bool, error) {
	switch s.config.Abi_driver {
	case "sysfs":
		{
			valid_line, err := cast.ToUint8E(line)
			if err != nil {
				return false, err
			}
			if !s.alredayOpen(line) {
				if err := s.exportBySysfs(valid_line); err != nil {
					return false, err
				}
				// for sysfs interface, leave nil object
				s.openedLine[line] = &gpiod.Line{}
				// waiting for gpio device fd
				time.Sleep(1 * time.Second)
			}
			if err = s.setDirectionBySysfs(valid_line, "in"); err != nil {
				return false, err
			}
			return s.getValueBySysfs(valid_line)
		}
	case "chardev":
		{
			valid_line, err := cast.ToIntE(line)
			if err != nil {
				return false, err
			}
			if !s.alredayOpen(line) {
				if s.openedChip != nil {
					return false, errors.New("using invalid chip")
				}
				l, err := s.openedChip.RequestLine(valid_line, gpiod.AsInput)
				if err != nil {
					return false, err
				}
				s.openedLine[line] = l
			}
			return s.getValueByChardev(s.openedLine[line])
		}
	}
	return false, errors.New("invalid interface")
}

func (s *Driver) setGPIO(line string, value bool) error {
	switch s.config.Abi_driver {
	case "sysfs":
		{
			valid_line, err := cast.ToUint8E(line)
			if err != nil {
				return err
			}
			if !s.alredayOpen(line) {
				if err := s.exportBySysfs(valid_line); err != nil {
					return err
				}
			}
			if err = s.setDirectionBySysfs(valid_line, "out"); err != nil {
				return err
			}
			return s.setValueBySysfs(valid_line, value)
		}
	case "chardev":
		{
			valid_line, err := cast.ToIntE(line)
			if err != nil {
				return err
			}
			ctx := 1
			if !s.alredayOpen(line) {
				if s.openedChip != nil {
					return errors.New("using invalid chip")
				}
				if !value {
					ctx = 0
				}
				l, err := s.openedChip.RequestLine(valid_line, gpiod.AsOutput(ctx))
				if err != nil {
					return err
				}
				s.openedLine[line] = l
			}
			return s.setValueByChardev(s.openedLine[line], ctx)
		}
	}
	return errors.New("invalid interface")
}

func (s *Driver) alredayOpen(line string) bool {
	for l := range s.openedLine {
		if line == l {
			return true
		}
	}
	return false
}
