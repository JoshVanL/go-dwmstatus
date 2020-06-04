package memory

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
)

func Memory(h *handler.Handler, s *string) error {
	update := func() {
		mem := h.SysInfo().Memory()
		*s = fmt.Sprintf(" %.1f/%.1f",
			mem[0], mem[1])
	}

	//ticker := time.NewTicker(time.Second / 2)
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			<-ticker.C
			update()
			h.Tick()
		}
	}()

	return nil
}
