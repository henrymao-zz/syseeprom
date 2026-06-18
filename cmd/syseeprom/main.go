package main

import (
	"fmt"
	"log"
	"os"

	"platform/eeprom"
	"platform/sysconf"
	_ "platform/vendors/broadcom"
	_ "platform/vendors/dell/s5232f"
	_ "platform/vendors/mellanox"
)

func main() {
	detectedVendor := resolveVendor()

	i2cPath := "/sys/bus/i2c/devices/i2c-0/0-0050/eeprom"

	driver, err := eeprom.GetDriver(detectedVendor, i2cPath)
	if err != nil {
		log.Fatalf("Failed to load EEPROM driver: %v", err)
	}

	mac, err := driver.GetBaseMAC()
	if err != nil {
		fmt.Printf("Base MAC (unavailable): %v\n", err)
	} else {
		fmt.Printf("System Base MAC: %s\n", mac)
	}

	sn, err := driver.GetSerialNumber()
	if err != nil {
		fmt.Printf("Serial Number (unavailable): %v\n", err)
	} else {
		fmt.Printf("Serial Number: %s\n", sn)
	}

	model, err := driver.GetModel()
	if err != nil {
		fmt.Printf("Model (unavailable): %v\n", err)
	} else {
		fmt.Printf("Model: %s\n", model)
	}
}

func resolveVendor() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	conf, err := sysconf.ReadMachineConf("/etc/machine.conf")
	if err != nil {
		log.Printf("Cannot read /etc/machine.conf: %v", err)
		return "mellanox"
	}

	onieMachine, ok := conf["onie_machine"]
	if !ok {
		log.Print("onie_machine not found in /etc/machine.conf")
		return "mellanox"
	}

	driver, err := sysconf.ResolveDriver(onieMachine)
	if err != nil {
		log.Printf("Cannot resolve driver for onie_machine %q: %v", onieMachine, err)
		return "mellanox"
	}

	log.Printf("Detected onie_machine=%s → driver=%s", onieMachine, driver)
	return driver
}
