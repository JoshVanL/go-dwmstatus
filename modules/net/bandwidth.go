package net

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/utils"
)

const (
	devPath     = "/sys/class/net"
	statsRXPath = "statistics/rx_bytes"
	statsTXPath = "statistics/tx_bytes"

	secs = 2
)

func Bandwidth(h *handler.Handler, s *string) error {
	ifaceHandler, err := getNetHandler()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second * secs)

	// received, transmitted
	x := getBytesX(ifaceHandler)

	*s = "🔻0KiB 🔺0KiB"
	h.Tick()

	go func() {
		for {
			currX := getBytesX(ifaceHandler)

			if !(currX[0] == x[0] && currX[1] == x[1]) {
				*s = fmt.Sprintf("🔻%.1fKiB 🔺%.1fKiB", (currX[0]-x[0])/(1024*2), (currX[1]-x[1])/(1024*2))
				x[0], x[1] = currX[0], currX[1]

				h.Tick()
			}

			<-ticker.C
		}
	}()

	return nil
}

func getBytesX(ifaceHandler *netHandler) [2]float64 {
	var currX [2]float64
	for _, iface := range ifaceHandler.ifaces {
		for i, xf := range []string{
			statsRXPath, statsTXPath,
		} {

			fpath := filepath.Join(devPath, iface.Name, xf)
			x, err := utils.ReadFileToFloat(fpath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read %s: %s", fpath, err)
				continue
			}

			currX[i] += x
		}
	}

	return currX
}
