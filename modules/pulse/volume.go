package pulse

import (
	"fmt"
	"math"
	"os"

	"github.com/joshvanl/go-dwmstatus/handler"
)

func Volume(h *handler.Handler, s *string) error {
	handler, err := getHandler()
	if err != nil {
		return err
	}

	go func() {
		for {
			*s, err = updateVolume(handler)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update volume: %s\n", err)
			}

			h.Tick()
			<-handler.events
		}
	}()

	return nil
}

func updateVolume(handler *pulseHandler) (string, error) {
	sinks, err := handler.GetSinks()
	if err != nil {
		return "-%", err
	}

	if len(sinks) == 0 {
		return "", nil
	}

	lastSink := sinks[len(sinks)-1]

	if lastSink.Muted {
		return " x", nil
	}

	vol := math.Round(100 * float64(lastSink.CurrentVolumeStep) / float64(lastSink.NumVolumeSteps))

	icon := ""
	if vol == 0 {
		icon = ""
	} else if vol < 40 {
		icon = ""
	}

	return fmt.Sprintf(" %s %.0f%%", icon, vol), nil
}
