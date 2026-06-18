package eeprom

type EEPROM interface {
	GetBaseMAC() (string, error)
	GetSerialNumber() (string, error)
	GetModel() (string, error)
}
