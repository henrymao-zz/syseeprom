package broadcom

import (
	"fmt"

	"platform/eeprom"
)

type Driver struct {
	i2cBus string
}

func NewDriver(bus string) *Driver {
	return &Driver{i2cBus: bus}
}

func (b *Driver) GetBaseMAC() (string, error) {
	data, err := readONIEEEPROM(b.i2cBus, 0x20, 6)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		data[0], data[1], data[2], data[3], data[4], data[5]), nil
}

func (b *Driver) GetSerialNumber() (string, error) {
	data, err := readONIEEEPROM(b.i2cBus, 0x2E, 16)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *Driver) GetModel() (string, error) {
	data, err := readONIEEEPROM(b.i2cBus, 0x3E, 16)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readONIEEEPROM(bus string, offset, length int) ([]byte, error) {
	_ = bus
	_ = offset
	_ = length
	return []byte{}, fmt.Errorf("Broadcom ONIE EEPROM not available")
}

func init() {
	eeprom.Register("broadcom", func(bus string) eeprom.EEPROM {
		return NewDriver(bus)
	})
}
