package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	switch runtime.GOOS {
	case "darwin":
		airportDir := "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources"
		airportExe := "airport"
		if _, err := os.Stat(filepath.Join(airportDir, airportExe)); os.IsNotExist(err) {
			fmt.Println("airport does not exist")
			return
		}

		airportOut, err := exec.Command(filepath.Join(airportDir, airportExe), "-I").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		var ssid string
		reader := bufio.NewReader(bytes.NewBuffer(airportOut))
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				fmt.Fprintln(os.Stderr, "Not found SSID")
				return
			}
			if strings.Contains(string(line), "SSID") && !strings.Contains(string(line), "BSSID") {
				ssid = strings.Split(string(line), ": ")[1]
				break
			}
		}

		password, err := exec.Command("security", "find-generic-password", "-l", ssid, "-D", "AirPort network password", "-w").Output()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		fmt.Print(string(password))
	default:
		fmt.Fprintln(os.Stderr, "Not support")
		return
	}
}
