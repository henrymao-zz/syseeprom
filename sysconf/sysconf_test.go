package sysconf

import (
	"os"
	"testing"
)

func TestReadMachineConf(t *testing.T) {
	content := `onie_arch=x86_64
onie_boot_reason=install
onie_machine=dellemc_s5232f_c3538
onie_platform=x86_64-dellemc_s5232f_c3538-r0
onie_version=3.40.1.1-9

# comment line
nos_name=Ubuntu
`

	tmp, err := os.CreateTemp("", "machine.conf.*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmp.Close()

	conf, err := ReadMachineConf(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]string{
		"onie_arch":     "x86_64",
		"onie_machine":  "dellemc_s5232f_c3538",
		"onie_platform": "x86_64-dellemc_s5232f_c3538-r0",
		"onie_version":  "3.40.1.1-9",
		"nos_name":      "Ubuntu",
	}

	for key, expected := range tests {
		if v := conf[key]; v != expected {
			t.Errorf("conf[%q] = %q, want %q", key, v, expected)
		}
	}
}

func TestResolveDriver(t *testing.T) {
	tests := []struct {
		onieMachine string
		expected    string
	}{
		{"dellemc_s5232f_c3538", "dell-s5232f"},
		{"dellemc_s5200_c3538", "dell-s5232f"},
	}

	for _, tt := range tests {
		driver, err := ResolveDriver(tt.onieMachine)
		if err != nil {
			t.Errorf("ResolveDriver(%q) error: %v", tt.onieMachine, err)
			continue
		}
		if driver != tt.expected {
			t.Errorf("ResolveDriver(%q) = %q, want %q", tt.onieMachine, driver, tt.expected)
		}
	}
}

func TestResolveDriverUnknown(t *testing.T) {
	_, err := ResolveDriver("unknown_vendor_xyz")
	if err == nil {
		t.Error("expected error for unknown onie_machine")
	}
}
