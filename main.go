package main

import (
	"fmt"

	"github.com/joshvanl/go-dwmstatus/errors"
	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/backlight"
	"github.com/joshvanl/go-dwmstatus/modules/battery"
	"github.com/joshvanl/go-dwmstatus/modules/bluetooth"
	"github.com/joshvanl/go-dwmstatus/modules/cpu"
	"github.com/joshvanl/go-dwmstatus/modules/datetime"
	"github.com/joshvanl/go-dwmstatus/modules/disk"
	"github.com/joshvanl/go-dwmstatus/modules/memory"
	"github.com/joshvanl/go-dwmstatus/modules/net"
	"github.com/joshvanl/go-dwmstatus/modules/temp"
	"github.com/joshvanl/go-dwmstatus/modules/volume"
	"github.com/joshvanl/go-dwmstatus/modules/weather"
)

var (
	enabledBlocks = []struct {
		f func(*handler.Handler, *string) error
		string
	}{
		{weather.Weather, "weather"},
		{sep, ""},
		{bluetooth.Bluetooth, "bluetooth"},
		{volume.Mic, "mic"},
		{volume.Volume, "volume"},
		{sep, ""},
		{cpu.CPU, "cpu"},
		{space, ""},
		{memory.Memory, "memory"},
		{space, ""},
		{disk.Disk, "disk"},
		{sep, ""},
		{net.Wifi, "wifi"},
		{net.Bandwidth, "bandwidth"},
		{sep, ""},
		{net.IFace, "iface"},
		{sep, ""},
		{temp.Temp, "temperature"},
		{sep, ""},
		{backlight.Backlight, "backlight"},
		{sep, ""},
		{battery.Battery, "battery"},
		{sep, ""},
		{datetime.DateTime, "datetime"},
	}
)

func sep(_ *handler.Handler, s *string) error {
	*s = " | "
	return nil
}

func space(_ *handler.Handler, s *string) error {
	*s = " "
	return nil
}

func main() {
	h, err := handler.New()
	if err != nil {
		errors.Kill(fmt.Errorf("error creating handler: %s\n", err))
	}

	for _, registerModule := range enabledBlocks {
		if len(registerModule.string) > 0 {
			fmt.Printf("registering module: %s\n", registerModule.string)
		}

		s := h.NewModule()
		h.Must(registerModule.f(h, s))
	}

	fmt.Printf("all modules registered\n")

	select {}
}
