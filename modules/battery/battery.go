package battery

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/utils"
)

const (
	path        = "/sys/class/power_supply"
	batteryName = "BAT0"
)

var (
	capPath  = filepath.Join(path, batteryName, "capacity")
	statPath = filepath.Join(path, batteryName, "status")
)

type battery struct {
	s *string
	h *handler.Handler
}

// TODO: make better..
func Battery(h *handler.Handler, s *string) error {
	ticker := time.NewTicker(time.Second * 2)

	b := &battery{
		s: s,
		h: h,
	}

	if err := b.setBatteryString(); err != nil {
		return err
	}
	b.h.Tick()

	go func() {
		for {
			<-ticker.C
			b.h.Must(b.setBatteryString())
			b.h.Tick()
		}
	}()

	return nil
}

func (b *battery) setBatteryString() error {
	status, capacity := b.getFiles()
	i, err := strconv.Atoi(string(capacity))
	if err != nil {
		return err
	}

	bat := getIcon(i)
	var charging string
	if string(status) == "Charging" {
		charging = " "
	}

	*b.s = fmt.Sprintf("%s%s %s%%", bat, charging, capacity)

	return nil
}

func (b *battery) getFiles() (status, capacity []byte) {
	status, err := utils.ReadFile(statPath)
	b.h.Must(err)

	capacity, err = utils.ReadFile(capPath)
	b.h.Must(err)

	if string(capacity) == "100" {
		status = []byte("full")
	}

	return status, capacity
}

func getIcon(capacity int) string {
	switch {
	case capacity == 100:
		//b.Color = "#FFFFFF"
		//b.Color = "#000000"
		return ""

	//case capacity > 90:
	//b.Color = "#ccffcc"

	//case capacity > 70:
	//b.Color = "#bbffbb"

	case capacity > 75:
		return ""

	case capacity > 60:
		//b.Color = "#ddffaa"

	case capacity > 50:
		//b.Color = "#eeffaa"

	case capacity > 30:
		//b.Color = "#ffdd77"
		return ""

	case capacity >= 25:
		//b.Color = "#ffaaaa"
		return ""
	}

	//b.Color = "#FF0000"
	return ""
}
