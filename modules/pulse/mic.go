package pulse

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joshvanl/go-dwmstatus/handler"
)

func Mic(h *handler.Handler, s *string) error {
	handler, err := getHandler()
	if err != nil {
		return err
	}

	go func() {
		for {
			*s, err = updateMic(handler)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update microphone: %s\n", err)
			}
			h.Tick()

			<-handler.events
		}
	}()

	return nil
}

func updateMic(handler *pulseHandler) (string, error) {
	sinks, err := handler.GetSources()
	if err != nil {
		return "-", err
	}

	for _, s := range sinks {
		if strings.Contains(s.Name, "alsa_input") {
			if s.Muted {
				return "", nil
			}
			return "", nil
		}
	}

	return "-", errors.New("failed to find microphone")
}
