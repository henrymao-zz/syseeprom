package s5232f

import (
	"encoding/binary"
	"fmt"
	"os"

	"platform/eeprom"
)

const (
	tlvCodeProductName   uint8 = 0x21
	tlvCodePartNumber    uint8 = 0x22
	tlvCodeSerialNumber  uint8 = 0x23
	tlvCodeMACBase       uint8 = 0x24
	tlvCodeDeviceVersion uint8 = 0x26
	tlvCodeServiceTag    uint8 = 0x2F
	tlvCodeVendorExt     uint8 = 0xFD
	tlvCodeCRC32         uint8 = 0xFE

	tlvInfoHdrLen     = 11
	tlvInfoMaxLen     = 2048
	tlvInfoIDString   = "TlvInfo\x00"
	tlvInfoVersion    = 0x01
	tlvMinLen         = 2
)

var idStringBytes = []byte(tlvInfoIDString)

type Driver struct {
	eepromPath string
	rawData    []byte
	tlvMap     map[uint8][]byte
}

func NewDriver(bus string) *Driver {
	d := &Driver{}
	if err := d.probeEEPROM(bus); err != nil {
		return d
	}
	if err := d.readAndParse(); err != nil {
		return d
	}
	return d
}

func (d *Driver) probeEEPROM(bus string) error {
	paths := buildProbePaths(bus)
	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			d.eepromPath = p
			return nil
		}
	}
	return fmt.Errorf("Dell S5232f EEPROM not found on any probed I2C bus")
}

func buildProbePaths(bus string) []string {
	var paths []string
	if bus != "" {
		paths = append(paths, bus)
	}
	for _, b := range []int{0, 1} {
		p := fmt.Sprintf("/sys/bus/i2c/devices/i2c-%d/%d-0050/eeprom", b, b)
		if p != bus {
			paths = append(paths, p)
		}
	}
	return paths
}

func (d *Driver) readAndParse() error {
	if d.eepromPath == "" {
		return fmt.Errorf("no EEPROM path configured")
	}

	data, err := os.ReadFile(d.eepromPath)
	if err != nil {
		return fmt.Errorf("failed to read EEPROM at %s: %w", d.eepromPath, err)
	}

	if len(data) < tlvInfoHdrLen {
		return fmt.Errorf("EEPROM data too short (%d bytes)", len(data))
	}

	if !isValidTLVHeader(data) {
		return fmt.Errorf("invalid TlvInfo header")
	}

	d.rawData = data
	d.parseTLVs()
	return nil
}

func isValidTLVHeader(data []byte) bool {
	if len(data) < tlvInfoHdrLen {
		return false
	}
	for i := 0; i < 8; i++ {
		if data[i] != idStringBytes[i] {
			return false
		}
	}
	if data[8] != tlvInfoVersion {
		return false
	}
	totalLen := binary.BigEndian.Uint16(data[9:11])
	return int(totalLen) <= tlvInfoMaxLen-tlvInfoHdrLen
}

func (d *Driver) parseTLVs() {
	d.tlvMap = make(map[uint8][]byte)

	totalLen := int(binary.BigEndian.Uint16(d.rawData[9:11]))
	tlvEnd := tlvInfoHdrLen + totalLen
	dataLen := len(d.rawData)

	tlvIndex := tlvInfoHdrLen
	for tlvIndex+tlvMinLen <= dataLen && tlvIndex < tlvEnd {
		code := d.rawData[tlvIndex]
		length := int(d.rawData[tlvIndex+1])

		if tlvIndex+2+length > dataLen {
			break
		}

		value := make([]byte, length)
		copy(value, d.rawData[tlvIndex+2:tlvIndex+2+length])
		d.tlvMap[code] = value

		if code == tlvCodeCRC32 {
			break
		}

		tlvIndex += length + 2
	}
}

func (d *Driver) getTLVField(code uint8) ([]byte, error) {
	if d.tlvMap == nil {
		return nil, fmt.Errorf("EEPROM not parsed")
	}
	value, exists := d.tlvMap[code]
	if !exists {
		return nil, fmt.Errorf("TLV field 0x%02X not found", code)
	}
	return value, nil
}

func (d *Driver) GetBaseMAC() (string, error) {
	value, err := d.getTLVField(tlvCodeMACBase)
	if err != nil {
		return "", err
	}
	if len(value) < 6 {
		return "", fmt.Errorf("MAC address TLV too short: %d bytes", len(value))
	}
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X",
		value[0], value[1], value[2], value[3], value[4], value[5]), nil
}

func (d *Driver) GetSerialNumber() (string, error) {
	value, err := d.getTLVField(tlvCodeSerialNumber)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (d *Driver) GetModel() (string, error) {
	value, err := d.getTLVField(tlvCodeProductName)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (d *Driver) GetPartNumber() (string, error) {
	value, err := d.getTLVField(tlvCodePartNumber)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (d *Driver) GetServiceTag() (string, error) {
	value, err := d.getTLVField(tlvCodeServiceTag)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (d *Driver) GetDeviceVersion() (string, error) {
	value, err := d.getTLVField(tlvCodeDeviceVersion)
	if err != nil {
		return "", err
	}
	if len(value) == 0 {
		return "", fmt.Errorf("Device Version TLV is empty")
	}
	return fmt.Sprintf("%d", value[0]), nil
}

func (d *Driver) SystemEEPROMInfo() map[string]string {
	result := make(map[string]string, len(d.tlvMap))
	for code, value := range d.tlvMap {
		result[fmt.Sprintf("0x%02X", code)] = string(value)
	}
	return result
}

func (d *Driver) EEPROMPath() string {
	return d.eepromPath
}

func init() {
	eeprom.Register("dell-s5232f", func(bus string) eeprom.EEPROM {
		return NewDriver(bus)
	})
}
