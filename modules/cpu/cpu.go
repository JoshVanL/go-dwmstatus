package cpu

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
)

func CPU(h *handler.Handler, s *string) error {
	update := func() {
		load := h.SysInfo().CPULoad()

		if load == -1 {
			*s = " 0.00%"
		} else {
			*s = fmt.Sprintf(" %.2f%%", load)
		}
	}

	//ticker := time.NewTicker(time.Second / 2)
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			update()
			h.Tick()
			<-ticker.C
		}
	}()

	return nil
}
