package bandwidth

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
	statPath  = "/proc/net/dev"
	secs      = 2
)

type netStat struct {
	h *handler.Handler

	down        bool
	received    uint64
	transmitted uint64
}

func Bandwidth(h *handler.Handler, s *string) error {
	ch := h.WatchSignal(watcher.RealTimeSignals["RTMIN+1"])

	n := netStat{
		h:    h,
		down: ifaceDown(h),
	}

	ticker := time.NewTicker(secs)

	go func() {
		for {
			n.down = ifaceDown(n.h)
			n.setString(s)
			h.Tick()

			select {
			case <-ch:
			case <-ticker.C:
			}
		}
	}()

	return nil
}

func ifaceDown(h *handler.Handler) bool {
	b, err := utils.ReadFile(filepath.Join(
		"/sys/class/net", ifaceName, "operstate"))
	h.Must(err)

	return string(b) == "down"
}

func (n *netStat) setString(s *string) {
	if n.down {
		*s = ""
		return
	}

	b, err := utils.ReadFile(statPath)
	n.h.Must(err)

	for _, lines := range strings.Split(string(b), "\n") {
		fields := strings.Fields(lines)

		if len(fields) < 4 || fields[0] != ifaceName+":" {
			continue
		}

		rec, err := strconv.ParseUint(fields[1], 10, 64)
		n.h.Must(err)

		tran, err := strconv.ParseUint(fields[9], 10, 64)
		n.h.Must(err)

		recf := float64(rec-n.received) / (secs * 1024 * 1024)
		tranf := float64(tran-n.transmitted) / (secs * 1024 * 1024)

		if n.received != 0 {
			*s = fmt.Sprintf(" %.1f / %.1f Mb/s",
				recf, tranf)
		} else {
			*s = " 0.0 / 0.0 Mb/s"
		}

		n.received = rec
		n.transmitted = tran

		return
	}
}
