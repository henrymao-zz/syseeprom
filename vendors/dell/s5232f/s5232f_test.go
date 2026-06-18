package s5232f

import (
	"encoding/binary"
	"testing"
)

func buildTestEEPROM(tlvs map[uint8][]byte) []byte {
	buf := make([]byte, 0, tlvInfoHdrLen+256)

	buf = append(buf, idStringBytes...)
	buf = append(buf, tlvInfoVersion)

	tlvData := make([]byte, 0, 256)
	for code, value := range tlvs {
		tlvData = append(tlvData, code)
		tlvData = append(tlvData, uint8(len(value)))
		tlvData = append(tlvData, value...)
	}

	totalLen := len(tlvData)
	lenBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(lenBytes, uint16(totalLen))
	buf = append(buf, lenBytes...)
	buf = append(buf, tlvData...)

	return buf
}

func TestParseTLVs(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeProductName:  []byte("S5232F-ON"),
		tlvCodePartNumber:   []byte("0J9K3T"),
		tlvCodeSerialNumber: []byte("CN0J9K3TABC001"),
		tlvCodeMACBase:      {0x00, 0xE0, 0xEC, 0x12, 0x34, 0x56},
		tlvCodeServiceTag:   []byte("ABC123D"),
		tlvCodeDeviceVersion: {0x04},
	}

	eepromData := buildTestEEPROM(tlvs)

	d := &Driver{}
	d.rawData = eepromData
	d.eepromPath = "/fake/path"

	if !isValidTLVHeader(d.rawData) {
		t.Fatal("valid TLV header rejected")
	}

	d.parseTLVs()

	if len(d.tlvMap) != len(tlvs) {
		t.Fatalf("expected %d TLVs, got %d", len(tlvs), len(d.tlvMap))
	}
}

func TestGetBaseMAC(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeMACBase:     {0x14, 0x18, 0x77, 0xAB, 0xCD, 0xEF},
		tlvCodeProductName: []byte("S5232F-ON"),
		tlvCodeSerialNumber: []byte("TEST123"),
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	mac, err := d.GetBaseMAC()
	if err != nil {
		t.Fatalf("GetBaseMAC failed: %v", err)
	}
	expected := "14:18:77:AB:CD:EF"
	if mac != expected {
		t.Fatalf("expected %q, got %q", expected, mac)
	}
}

func TestGetModel(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeProductName:  []byte("S5232F-ON"),
		tlvCodeMACBase:      {0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		tlvCodeSerialNumber: []byte("TEST123"),
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	model, err := d.GetModel()
	if err != nil {
		t.Fatalf("GetModel failed: %v", err)
	}
	if model != "S5232F-ON" {
		t.Fatalf("expected %q, got %q", "S5232F-ON", model)
	}
}

func TestGetSerialNumber(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeSerialNumber: []byte("CN0J9K3TABC001"),
		tlvCodeProductName:  []byte("S5232F-ON"),
		tlvCodeMACBase:      {0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	sn, err := d.GetSerialNumber()
	if err != nil {
		t.Fatalf("GetSerialNumber failed: %v", err)
	}
	if sn != "CN0J9K3TABC001" {
		t.Fatalf("expected %q, got %q", "CN0J9K3TABC001", sn)
	}
}

func TestGetPartNumber(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodePartNumber:   []byte("0J9K3T"),
		tlvCodeProductName:  []byte("S5232F-ON"),
		tlvCodeMACBase:      {0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		tlvCodeSerialNumber: []byte("TEST123"),
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	pn, err := d.GetPartNumber()
	if err != nil {
		t.Fatalf("GetPartNumber failed: %v", err)
	}
	if pn != "0J9K3T" {
		t.Fatalf("expected %q, got %q", "0J9K3T", pn)
	}
}

func TestGetServiceTag(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeServiceTag:   []byte("ABC123D"),
		tlvCodeProductName:  []byte("S5232F-ON"),
		tlvCodeMACBase:      {0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		tlvCodeSerialNumber: []byte("TEST123"),
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	st, err := d.GetServiceTag()
	if err != nil {
		t.Fatalf("GetServiceTag failed: %v", err)
	}
	if st != "ABC123D" {
		t.Fatalf("expected %q, got %q", "ABC123D", st)
	}
}

func TestGetDeviceVersion(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeDeviceVersion: {0x04},
		tlvCodeProductName:   []byte("S5232F-ON"),
		tlvCodeMACBase:       {0x00, 0x01, 0x02, 0x03, 0x04, 0x05},
		tlvCodeSerialNumber:  []byte("TEST123"),
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	dv, err := d.GetDeviceVersion()
	if err != nil {
		t.Fatalf("GetDeviceVersion failed: %v", err)
	}
	if dv != "4" {
		t.Fatalf("expected %q, got %q", "4", dv)
	}
}

func TestIsValidTLVHeader(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name: "valid header",
			data: buildTestEEPROM(map[uint8][]byte{
				tlvCodeProductName: []byte("test"),
			}),
			expected: true,
		},
		{
			name:     "too short",
			data:     []byte{0x00, 0x01},
			expected: false,
		},
		{
			name: "wrong ID string",
			data: func() []byte {
				data := make([]byte, tlvInfoHdrLen)
				copy(data, []byte("BadData\x00"))
				data[8] = tlvInfoVersion
				return data
			}(),
			expected: false,
		},
		{
			name: "wrong version",
			data: func() []byte {
				data := make([]byte, tlvInfoHdrLen)
				copy(data, idStringBytes)
				data[8] = 0xFF
				return data
			}(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidTLVHeader(tt.data)
			if result != tt.expected {
				t.Errorf("isValidTLVHeader() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSystemEEPROMInfo(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeProductName:  []byte("S5232F-ON"),
		tlvCodePartNumber:   []byte("0J9K3T"),
		tlvCodeSerialNumber: []byte("CN0J9K3TABC001"),
		tlvCodeMACBase:      {0x00, 0xE0, 0xEC, 0x12, 0x34, 0x56},
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	info := d.SystemEEPROMInfo()
	if len(info) != len(tlvs) {
		t.Fatalf("expected %d entries in info map, got %d", len(tlvs), len(info))
	}

	if info["0x21"] != "S5232F-ON" {
		t.Errorf("expected product name 'S5232F-ON', got %q", info["0x21"])
	}
}

func TestMissingTLVField(t *testing.T) {
	tlvs := map[uint8][]byte{
		tlvCodeProductName: []byte("S5232F-ON"),
	}

	d := &Driver{
		rawData:    buildTestEEPROM(tlvs),
		eepromPath: "/fake/path",
	}
	d.parseTLVs()

	_, err := d.GetSerialNumber()
	if err == nil {
		t.Error("expected error for missing serial number TLV")
	}
}
