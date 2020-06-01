package volume

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/auroralaboratories/pulse"
	"github.com/joshvanl/go-dwmstatus/handler"
)

func Mic(h *handler.Handler, s *string) error {
	sinkClient, err := pulse.NewClient("go-dwnstatus-mic")
	if err != nil {
		return fmt.Errorf("failed to get mic client for pulse: %s", err)
	}

	ch, err := pulseWatcher()
	if err != nil {
		return fmt.Errorf("failed to get watcher for mic: %s", err)
	}

	go func() {
		for {
			*s, err = updateMic(sinkClient)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update microphone: %s", err)
			}
			h.Tick()

			<-ch
		}
	}()

	return nil
}

func updateMic(c *pulse.Client) (string, error) {
	sinks, err := c.GetSources()
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
