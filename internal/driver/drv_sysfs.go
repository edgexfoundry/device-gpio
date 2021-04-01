// +build sysfs

package driver

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/spf13/cast"
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
	valid_port, err := cast.ToUint8E(split[1])
	if err != nil {
		return err
	}
	if err := export(valid_port); err != nil {
		return err
	}
	if err = setDirection(valid_port, "out"); err != nil {
		return err
	}
	return setValue(valid_port, value)
}

func (dev *GPIODev) GetGPIO(gpio string) (bool, error) {
	if !strings.Contains(gpio, "/") {
		return false, errors.New("invalid gpio number")
	}
	split := strings.Split(gpio, "/")
	valid_port, err := cast.ToUint8E(split[1])
	if err != nil {
		return false, err
	}
	if err := export(valid_port); err != nil {
		return false, err
	}
	if err = setDirection(valid_port, "in"); err != nil {
		return false, err
	}
	return getValue(valid_port)
}

func export(gpioNum uint8) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	return ioutil.WriteFile("/sys/class/gpio/export", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
}

func unexport(gpioNum uint8) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return ioutil.WriteFile("/sys/class/gpio/unexport", []byte(fmt.Sprintf("%d\n", gpioNum)), 0644)
	}
	return nil
}

func setDirection(gpioNum uint8, direction string) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
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
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum), []byte(way), 0644)
	} else {
		return errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func getDirection(gpioNum uint8) (string, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		direction, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/direction", gpioNum))
		if err != nil {
			return "", err
		} else {
			return strings.Replace(string(direction), "\n", "", -1), err
		}
	} else {
		return "", errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func setValue(gpioNum uint8, value bool) error {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		var tmp string
		if value {
			tmp = "1"
		} else {
			tmp = "0"
		}
		return ioutil.WriteFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum), []byte(tmp), 0644)
	} else {
		return errors.New("unexpected behavior, the GPIO pin has not been exported")
	}
}

func getValue(gpioNum uint8) (bool, error) {
	path := fmt.Sprintf("/sys/class/gpio/gpio%d", gpioNum)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		ret, err := ioutil.ReadFile(fmt.Sprintf("/sys/class/gpio/gpio%d/value", gpioNum))
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
