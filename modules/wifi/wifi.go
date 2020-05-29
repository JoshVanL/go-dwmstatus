package wifi

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/utils"
	"github.com/joshvanl/go-dwmstatus/watcher"
)

const (
	//ifaceName = "wlan0"
	ifaceName = "enp0s20f0u8"
)

func Wifi(h *handler.Handler, s *string) error {
	update := func() error {
		n, err := wirelessInt(h)
		if err != nil {
			return err
		}

		if n == -1 {
			*s = ""
			return nil
		}

		//if n >= 70 {
		//	block.Color = "#0050b8"
		//} else if n >= 50 {
		if n >= 50 {
			//block.Color = "#aaddaa"
		} else if n >= 30 {
			//block.Color = "#a65f3d"
			//block.Color = "#ffae88"
		} else {
			//block.Color = "#ff0000"
		}

		//block.FullText = fmt.Sprintf(" %d%%", n)
		*s = fmt.Sprintf(" %d%%", n)
		return nil
	}

	if err := update(); err != nil {
		return err
	}
	h.Tick()

	//h.Scheduler().Register(time.Minute, update)
	ticker := time.NewTicker(time.Second * 10)

	ch := h.WatchSignal(watcher.RealTimeSignals["RTMIN+1"])

	go func() {
		for {
			update()
			h.Tick()

			select {
			case <-ch:
			case <-ticker.C:
			}
		}
	}()

	return nil
}

func wirelessInt(h *handler.Handler) (int, error) {
	b, err := utils.ReadFile(filepath.Join(
		"/sys/class/net", ifaceName, "operstate"))
	if err != nil {
		return -1, err
	}

	if string(b) == "down" {
		return -1, nil
	}

	b, err = utils.ReadFile("/proc/net/wireless")
	if err != nil {
		return -1, err
	}

	for _, line := range strings.Split(string(b), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 || fields[0] != ifaceName+":" {
			continue
		}

		n, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return -1, err
		}

		return int(n * 100 / 70), nil
	}

	return -1, nil
}
