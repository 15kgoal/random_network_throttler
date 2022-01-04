package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

var (
	DownloadCap = "22256"       // In kpbs
	UploadCap   = "99999999999" // In kpbs
	clear       = false
)

func init() {
	flag.BoolVar(&clear, "clear", false, "If present, clears the current wifi interface")
	flag.Parse()
}

func main() {
	deviceListR, err := exec.Command("nmcli", "device", "status").Output()
	if err != nil {
		fmt.Printf("error on mcli %s\n", err)
	}
	deviceList := string(deviceListR)
	splitLines := strings.Split(deviceList, "\n")
	var networkInterface string

	for _, l := range splitLines {
		connected, _ := regexp.Match(`\ connected\ `, []byte(l))
		wifi, _ := regexp.Match(`\ wifi\ `, []byte(l))
		if connected && wifi {
			re := regexp.MustCompile(`[^\ ]*`)
			networkInterface = string(re.Find([]byte(l)))
		}
	}

	if networkInterface == "" {
		fmt.Println("No connected networks found.")
	}

	var wonderShaperR []byte
	if clear {
		clearInterface(networkInterface)
	} else {
		rand.Seed(time.Now().UnixNano())
		randInt := rand.Intn(5)
		if randInt == 0 {
			fmt.Printf("Capping %s with wondershaper. random int=%d\n", networkInterface, randInt)
			wonderShaperR, err = exec.Command("wondershaper", networkInterface, DownloadCap, UploadCap).Output()
			if err != nil {
				fmt.Printf("error on wondershaper %s\n", err)
			}

			fmt.Println(string(wonderShaperR))
		} else {
			fmt.Printf("No Capping %s with wondershaper, performing clear instead. random int=%d\n", networkInterface, randInt)
			clearInterface(networkInterface)
		}
	}

}

func clearInterface(networkInterface string) {
	fmt.Printf("Clearing %s with wondershaper\n", networkInterface)
	wonderShaperR, err := exec.Command("wondershaper", "clear", networkInterface).Output()
	if err != nil {
		fmt.Printf("error on wondershaper clear %s\n", err)
	}

	fmt.Println(string(wonderShaperR))
}
