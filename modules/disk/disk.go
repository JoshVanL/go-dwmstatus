package disk

import (
	"fmt"
	"syscall"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
)

func Disk(h *handler.Handler, s *string) error {
	update := func() error {
		var stat syscall.Statfs_t
		if err := syscall.Statfs("/", &stat); err != nil {
			return err
		}

		*s = fmt.Sprintf(" %.2fG",
			float64(stat.Bavail*uint64(stat.Bsize))/(1024*1024*1024))

		return nil
	}

	if err := update(); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Minute * 10)

	go func() {
		for {
			<-ticker.C
			h.Must(update())
			h.Tick()
		}
	}()

	return nil
}
