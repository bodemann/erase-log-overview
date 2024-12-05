package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Device struct {
	Manufacturer         string
	Computer             string
	ComputerSerialNumber string
	Drive                string
	DriveSerialNumber    string
}

func main() {
	files, err := os.ReadDir("logs")
	if err != nil {
		panic(err)
	}
	var devices []Device
	for _, file := range files {
		devices = append(devices, parseEraseLog(file.Name()))
	}
	sort.Slice(devices, func(i, j int) bool {
		return devices[i].Computer < devices[j].Computer
	})
	printDevicesTable(devices)
}

func printDevicesTable(devices []Device) {
	fmt.Println("Alle Geräte gelöscht am 05.12.2024 von 10:00 bis 14:30 Uhr\n\n")
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Manufacturer", "Computer", "Computer-Serial-Number", "Drive", "Drive-Serial-Number"})
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, device := range devices {
		table.Append([]string{device.Manufacturer, device.Computer, device.ComputerSerialNumber, device.Drive, device.DriveSerialNumber})
	}
	table.Render()
}

func parseEraseLog(fileName string) Device {
	rawData, err := os.ReadFile(filepath.Join("logs", fileName))
	if err != nil {
		panic(err)
	}
	data := strings.Split(string(rawData), "\n")
	var device Device
	for _, line := range data {
		switch {
		case strings.Contains(line, "Manufacturer:"):
			device.Manufacturer = strings.TrimSpace(strings.Split(line, ":")[1])
		case strings.Contains(line, "Product Name:"):
			device.Computer = strings.TrimSpace(strings.Split(line, ":")[1])
		case strings.Contains(line, "Serial Number:"):
			device.ComputerSerialNumber = strings.TrimSpace(strings.Split(line, ":")[1])
		case strings.Contains(line, "/dev/nvme") || strings.Contains(line, "/dev/sda"):
			device.Drive = strings.TrimSpace(strings.Split(line, "(")[0])
			device.DriveSerialNumber = strings.TrimSpace(strings.Replace(strings.Split(line, ":")[1], " SIZE", "", -1))
		}
	}
	return device
}
