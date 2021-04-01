// build chardev

package driver

import (
	"errors"
	"fmt"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/spf13/cast"
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

type GPIODev struct {
	lc   logger.LoggingClient
	gpio uint8
}

func NewGPIODev(lc logger.LoggingClient) *GPIODev {
	return &GPIODev{lc: lc, gpio: 0}
}

func (dev *GPIODev) SetGPIO(gpio string, value bool) error {
	if !strings.Contains(gpio, "/") {
		return errors.New("invalid gpio number")
	}
	split := strings.Split(gpio, "/")
	valid_chip, err := cast.ToUint8E(split[0])
	if err != nil {
		return err
	}
	valid_line, err := cast.ToIntE(split[1])
	if err != nil {
		return err
	}
	fmt.Printf("chip: %d, line: %d\n", valid_chip, valid_line)
	return setValue(valid_chip, valid_line, value)
}

func (dev *GPIODev) GetGPIO(gpio string) (bool, error) {
	if !strings.Contains(gpio, "/") {
		return false, errors.New("invalid gpio number")
	}
	split := strings.Split(gpio, "/")
	valid_chip, err := cast.ToUint8E(split[0])
	if err != nil {
		return false, err
	}
	valid_line, err := cast.ToIntE(split[1])
	if err != nil {
		return false, err
	}
	return getValue(valid_chip, valid_line)
}

func setValue(chip uint8, line int, value bool) error {
	chipName := fmt.Sprintf("gpiochip%d", chip)
	c, err := gpiod.NewChip(chipName)
	defer c.Close()
	if err != nil {
		return err
	}
	ctx := 1
	if !value {
		ctx = 0
	}
	l, err := c.RequestLine(rpi.GPIO17, gpiod.AsOutput(ctx))
	if err := l.SetValue(ctx); err != nil {
		return err
	}
	fmt.Println("set ctx")
	defer l.Close()
	if err != nil {
		return err
	}
	return nil
}

func getValue(chip uint8, line int) (bool, error) {
	chipName := fmt.Sprintf("gpiochip%d", chip)
	c, err := gpiod.NewChip(chipName)
	defer c.Close()
	l, err := c.RequestLine(line, gpiod.AsInput)
	defer l.Close()
	if err != nil {
		return false, err
	}
	val, err := l.Value()
	if err != nil {
		return false, err
	}
	ctx := true
	if val == 0 {
		ctx = false
	}
	return ctx, nil
}
