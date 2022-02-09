package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/mdlayher/wol"
)

const (
	ENV_PORT      = "WEBWOL_PORT"
	ENV_ASSETS    = "WEBWOL_ASSETS_DIR"
	ENV_TEMPLATES = "WEBWOL_TEMPLATES_DIR"
	ENV_CONFIG    = "WEBWOL_CONFIG"
	ENV_BASEURL   = "WEBWOL_BASEURL"
)

type WakeUp struct {
	Device string `json:"device"`
	Mac    string `json:"mac"`
	Ip     string `json:"ip"`
}

var (
	data      []WakeUp
	configDir = "config"
)

func main() {
	log.Println("Starting WebWOL.")
	if ps, ok := os.LookupEnv(ENV_PORT); ok {
		p, err := strconv.Atoi(ps)
		if err == nil {
			port = p
		}
	}
	if ts, ok := os.LookupEnv(ENV_TEMPLATES); ok {
		templateDir = ts
	}
	if as, ok := os.LookupEnv(ENV_ASSETS); ok {
		assetsDir = as
	}
	if cs, ok := os.LookupEnv(ENV_CONFIG); ok {
		configDir = cs
	}
	if bs, ok := os.LookupEnv(ENV_BASEURL); ok {
		baseURL = bs
	}

	log.Printf("Using port %d, templates %s and assets %s", port, templateDir, assetsDir)

	http.HandleFunc("/", handlerIndex)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%04d", port), nil))
}

func wolUdp(ip string, mac string, password []byte) error {
	c, err := wol.NewClient()
	if err != nil {
		return err
	}
	defer c.Close()

	m, err := net.ParseMAC(mac)
	if password != nil {
		return c.WakePassword(ip, m, password)
	}
	return c.Wake(ip, m)
}

func loadData() {
	f, err := os.Open(configDir + "/wakeup.json")
	if err != nil {
		data = make([]WakeUp, 0)
	} else {
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			data = make([]WakeUp, 0)
		} else {
			err = json.Unmarshal(b, &data)
			if err != nil {
				data = make([]WakeUp, 0)
			}
		}
	}
	log.Printf("Having %d wakeups in config", len(data))
}

func insertOrUpdateData(wu WakeUp, odevice string, scope string) error {
	if wu.Device == "" || wu.Mac == "" || wu.Ip == "" {
		return fmt.Errorf("device, mac and ip must be set")
	}
	if scope == "insert" && deviceExists(wu.Device) {
		return fmt.Errorf("device %s exists already", wu.Device)
	}
	if scope == "update" && !deviceExists(odevice) {
		return fmt.Errorf("device %s does not exist", wu.Device)
	}
	if scope == "insert" {
		data = append(data, wu)
	} else {
		replaceWakeUp(wu, odevice)
	}
	return nil
}

func replaceWakeUp(wu WakeUp, odevice string) {
	for i, w := range data {
		if w.Device == odevice {
			data[i] = wu
			return
		}
	}
}

func deleteItem(device string) error {
	for i, w := range data {
		if w.Device == device {
			data = append(data[:i], data[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("device %s not found for deleting", device)
}

func cloneItem(device string) error {
	for _, w := range data {
		if w.Device == device {
			wu := WakeUp{
				Device: w.Device + "-clone",
				Mac:    w.Mac,
				Ip:     w.Ip,
			}
			data = append(data, wu)
			return nil
		}
	}
	return fmt.Errorf("device %s not found for cloning", device)
}

func saveData() error {
	f, err := os.Create(configDir + "/wakeup.json")
	if err == nil {
		defer f.Close()
		b, err := json.Marshal(data)
		if err == nil {
			f.Write(b)
			return nil
		} else {
			return fmt.Errorf("error writing file: %v", err)
		}
	} else {
		return fmt.Errorf("error creating config file: %v", err)
	}
}

func deviceExists(device string) bool {
	for _, w := range data {
		if w.Device == device {
			return true
		}
	}
	return false
}

func wakeupData(device string) (WakeUp, bool) {
	for _, w := range data {
		if w.Device == device {
			return w, true
		}
	}
	return WakeUp{}, false
}
