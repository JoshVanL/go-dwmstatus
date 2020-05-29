package main

import (
	"fmt"

	"github.com/joshvanl/go-dwmstatus/errors"
	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/backlight"
	"github.com/joshvanl/go-dwmstatus/modules/bandwidth"
	"github.com/joshvanl/go-dwmstatus/modules/battery"
	"github.com/joshvanl/go-dwmstatus/modules/bluetooth"
	"github.com/joshvanl/go-dwmstatus/modules/cpu"
	"github.com/joshvanl/go-dwmstatus/modules/datetime"
	"github.com/joshvanl/go-dwmstatus/modules/disk"
	"github.com/joshvanl/go-dwmstatus/modules/iface"
	"github.com/joshvanl/go-dwmstatus/modules/memory"
	"github.com/joshvanl/go-dwmstatus/modules/temp"
	"github.com/joshvanl/go-dwmstatus/modules/volume"
	"github.com/joshvanl/go-dwmstatus/modules/weather"
	//"github.com/joshvanl/go-dwmstatus/modules/wallpaper"
	"github.com/joshvanl/go-dwmstatus/modules/wifi"
)

var (
	enabledBlocks = []struct {
		f func(*handler.Handler, *string) error
		string
	}{
		{bluetooth.Bluetooth, "bluetooth"},
		{sep, ""},
		{volume.Mic, "mic"},
		{space, ""},
		{weather.Weather, "weather"},
		{sep, ""},
		{volume.Volume, "volume"},
		{sep, ""},
		{cpu.CPU, "cpu"},
		{space, ""},
		{memory.Memory, "memory"},
		{space, ""},
		{disk.Disk, "disk"},
		{sep, ""},
		{wifi.Wifi, "wifi"},
		{bandwidth.Bandwidth, "bandwidth"},
		{space, ""},
		{iface.IFace, "iface"},
		{sep, ""},
		{temp.Temp, "temperature"},
		{sep, ""},
		{backlight.Backlight, "backlight"},
		{sep, ""},
		{battery.Battery, "battery"},
		{sep, ""},
		{datetime.DateTime, "datetime"},
		{space, ""},
	}
)

func sep(_ *handler.Handler, s *string) error {
	*s = " | "
	return nil
}

func space(_ *handler.Handler, s *string) error {
	*s = "  "
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

	select {}
}
