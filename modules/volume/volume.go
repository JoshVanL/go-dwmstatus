package volume

import (
	"fmt"
	"math"
	"os"

	//"os"

	"github.com/auroralaboratories/pulse"

	"github.com/joshvanl/go-dwmstatus/handler"
)

var (
	watcherCh <-chan pulse.EventType
)

func Volume(h *handler.Handler, s *string) error {
	sinkClient, err := pulse.NewClient("go-dwnstatus-volume")
	//_, err := pulse.NewClient("go-dwnstatus-volume")
	if err != nil {
		return err
	}

	//ch := make(chan struct{})
	ch, err := pulseWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			*s, err = updateVolume(sinkClient)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update volume: %s\n", err)
			}

			h.Tick()
			<-ch
		}
	}()

	return nil
}

func pulseWatcher() (<-chan pulse.EventType, error) {
	if watcherCh != nil {
		return watcherCh, nil
	}

	subscriptionClient, err := pulse.NewClient("go-dwnstatus-watcher")
	if err != nil {
		return nil, err
	}

	watcherCh = subscriptionClient.Subscribe(pulse.AllEvent)
	return watcherCh, nil
}

func updateVolume(c *pulse.Client) (string, error) {
	sinks, err := c.GetSinks()
	if err != nil {
		return "-%", err
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
