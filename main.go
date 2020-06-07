package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/backlight"
	"github.com/joshvanl/go-dwmstatus/modules/battery"
	"github.com/joshvanl/go-dwmstatus/modules/bluetooth"
	"github.com/joshvanl/go-dwmstatus/modules/cpu"
	"github.com/joshvanl/go-dwmstatus/modules/datetime"
	"github.com/joshvanl/go-dwmstatus/modules/disk"
	"github.com/joshvanl/go-dwmstatus/modules/memory"
	"github.com/joshvanl/go-dwmstatus/modules/net"
	"github.com/joshvanl/go-dwmstatus/modules/pulse"
	"github.com/joshvanl/go-dwmstatus/modules/temp"
	"github.com/joshvanl/go-dwmstatus/modules/weather"
)

type module struct {
	register func(*handler.Handler, *string) error
	name     string
}

var (
	enabledModules = []*module{
		{weather.Weather, "weather"},
		{sep, ""},
		{bluetooth.Bluetooth, "bluetooth"},
		{pulse.Mic, "mic"},
		{pulse.Volume, "volume"},
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
		handler.Kill(fmt.Errorf("error creating handler: %s\n", err))
	}

	type failedRegister struct {
		*module
		*string
	}

	var retry []failedRegister

	for _, module := range enabledModules {
		if len(module.name) > 0 {
			fmt.Printf("registering module: %s\n", module.name)
		}

		s := h.NewModule()
		if !registerModule(h, module, s) {
			retry = append(retry, failedRegister{module, s})
		}
	}

	for len(retry) > 0 {
		time.Sleep(time.Second)

		var nextRetry []failedRegister
		for _, r := range retry {
			fmt.Printf("retrying registering module: %s\n", r.name)

			if !registerModule(h, r.module, r.string) {
				nextRetry = append(nextRetry, failedRegister{r.module, r.string})
			}
		}

		retry = nextRetry
	}

	fmt.Fprint(os.Stdout, "all modules registered\n")

	select {}
}

func registerModule(h *handler.Handler, module *module, s *string) bool {
	if err := module.register(h, s); err != nil {
		fmt.Fprintf(os.Stderr, "failed to register module %s: %s\n", module.name, err)
		return false
	}

	return true
}
