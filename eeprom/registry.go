package eeprom

import "fmt"

type DriverFactory func(bus string) EEPROM

var registry = make(map[string]DriverFactory)

func Register(vendorName string, factory DriverFactory) {
	registry[vendorName] = factory
}

func GetDriver(vendorName, bus string) (EEPROM, error) {
	factory, exists := registry[vendorName]
	if !exists {
		return nil, fmt.Errorf("unsupported vendor: %s", vendorName)
	}
	return factory(bus), nil
}
