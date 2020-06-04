package backlight

import (
	"fmt"
	"os"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/utils"
)

const (
	maxBFile = "/sys/class/backlight/intel_backlight/max_brightness"
	bFile    = "/sys/class/backlight/intel_backlight/brightness"
)

func Backlight(h *handler.Handler, s *string) error {
	ch, err := h.Watcher().Add(bFile)
	if err != nil {
		return err
	}

	go func() {
		for {
			*s, err = update(h)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update backlight: %s\n", err)
			}
			h.Tick()

			<-ch
		}
	}()

	return nil
}

func update(h *handler.Handler) (string, error) {
	max, err := utils.ReadFileToFloat(maxBFile)
	if err != nil {
		return "-", err
	}

	b, err := utils.ReadFileToFloat(bFile)
	if err != nil {
		return "-", err
	}

	return fmt.Sprintf(" %.0f%%", b*100/max), nil
}
