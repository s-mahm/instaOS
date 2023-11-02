package flash

import (
	"flag"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var (
	debug = flag.Int("debug", 0, "libusb debug level (0..3)")
)

func IsValidFlashDevice(device string) error {
	deviceModelVerifier := "ID_USB_MODEL=Flash_Disk"
	cmd := "udevadm"
	args := []string{"info", "-q", "property", "-n", device}
	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		return fmt.Errorf("checking device %s: %s", device, err)
	}
	if !strings.Contains(string(out), deviceModelVerifier) {
		return fmt.Errorf("invalid flash device %s", device)
	} else {
		deviceTypeVerifier := "DEVTYPE=disk"
		if !strings.Contains(string(out), deviceTypeVerifier) {
			return fmt.Errorf("please provide the device identifier excluding it's partition (i.e. /dev/sdb not /dev/sdb1)")
		}
	}
	return nil
}

func GetFlashDeviceInfo(device string) (string, error) {
	dfPattern, err := regexp.Compile(fmt.Sprintf(`(?m)^%s\d\s+(\d+)\s+\d+\s+\d+.*%%\s(.*)$`, device))
	if err != nil {
		return "", fmt.Errorf("getting device %s info: %s", device, err)
	}
	cmd := "df"
	out, err := exec.Command(cmd).Output()
	if err != nil {
		return "", fmt.Errorf("calling df: %s", err)
	}
	matches := dfPattern.FindStringSubmatch(string(out))
	if len(matches) == 0 {
		return "", fmt.Errorf("device %s not found in df output", device)
	}
	capacity, _ := strconv.ParseFloat(matches[1], 64)
	capacity = math.Ceil((capacity/1048576)*10) / 10
	if capacity < 5.0 {
		return "", fmt.Errorf("device %s needs to be at least 5GB in capacity", device)
	} else {
		fmt.Printf("Flash device %s found with capacity %.1fGB\n", device, capacity)
	}
	return matches[2], nil
}

func GetFlashDevices() {

}
