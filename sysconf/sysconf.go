package sysconf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadMachineConf(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		conf[parts[0]] = parts[1]
	}
	return conf, scanner.Err()
}

var onieMachineDriverMap = map[string]string{
	"dellemc_s5232f_c3538": "dell-s5232f",
	"dellemc_s5200_c3538":  "dell-s5232f",
}

func ResolveDriver(onieMachine string) (string, error) {
	if driver, ok := onieMachineDriverMap[onieMachine]; ok {
		return driver, nil
	}
	for key, driver := range onieMachineDriverMap {
		if strings.HasPrefix(onieMachine, key) || strings.HasPrefix(key, onieMachine) {
			return driver, nil
		}
	}
	return "", fmt.Errorf("no driver found for onie_machine: %s", onieMachine)
}
