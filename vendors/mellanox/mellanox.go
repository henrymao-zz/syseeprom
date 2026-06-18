package mellanox

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

func (m *Driver) GetBaseMAC() (string, error) {
	data, err := readEEPROM(m.i2cBus, 0x20, 6)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		data[0], data[1], data[2], data[3], data[4], data[5]), nil
}

func (m *Driver) GetSerialNumber() (string, error) {
	data, err := readEEPROM(m.i2cBus, 0x30, 16)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (m *Driver) GetModel() (string, error) {
	data, err := readEEPROM(m.i2cBus, 0x40, 16)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func readEEPROM(bus string, offset, length int) ([]byte, error) {
	_ = bus    // Simulated I2C read
	_ = offset // Simulated I2C read
	_ = length // Simulated I2C read
	return []byte{}, fmt.Errorf("Mellanox I2C EEPROM not available")
}

func init() {
	eeprom.Register("mellanox", func(bus string) eeprom.EEPROM {
		return NewDriver(bus)
	})
}
