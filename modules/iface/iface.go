package iface

import (
	"net"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/watcher"
)

const (
	//ifaceName = "wlan0"
	ifaceName = "enp0s20f0u8"
)

func IFace(h *handler.Handler, s *string) error {
	ch := h.WatchSignal(watcher.RealTimeSignals["RTMIN+1"])

	if err := update(s); err != nil {
		return err
	}
	h.Tick()

	ticker := time.NewTicker(time.Second * 60)

	go func() {
		for {
			update(s)
			h.Tick()

			select {
			case <-ch:
			case <-ticker.C:
			}
		}
	}()

	return nil
}

func update(s *string) error {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	if len(addrs) == 0 {
		*s = "down"
		return nil
		//block.Color = "#ee9999"
	}

	var found bool
	for _, addr := range addrs {
		v, ok := addr.(*net.IPNet)
		if !ok || v.IP.To4() == nil {
			continue
		}

		found = true
		*s = v.IP.String()
		//block.Color = "#aaddaa"
		break
	}

	if !found {
		*s = "IPERROR"
	}

	return nil
}
