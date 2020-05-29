package volume

import (
	"fmt"
	"math"
	"os"

	"github.com/auroralaboratories/pulse"

	"github.com/joshvanl/go-dwmstatus/handler"
)

var (
	watcherCh <-chan pulse.EventType
)

func Volume(h *handler.Handler, s *string) error {
	//sinkClient, err := pulse.NewClient("go-dwnstatus-volume")
	_, err := pulse.NewClient("go-dwnstatus-volume")
	if err != nil {
		return err
	}

	ch, err := pulseWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			//*s, err = "  ", nil
			*s = ""
			//*s, err = updateVolume(sinkClient)
			h.Tick()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to update volume: %s", err)
			}
			//TODO: add muted option

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
	vol := math.Round(100 * float64(sinks[len(sinks)-1].CurrentVolumeStep) / float64(sinks[len(sinks)-1].NumVolumeSteps))

	print(vol)
	return "  ", nil
	//icon := "  "
	//if vol == 0 {
	//icon = ""
	//} else if vol < 40 {
	//	icon = ""
	//}

	//return fmt.Sprintf(" %s %.0f%%", icon, vol), nil
	//return fmt.Sprintf(" %.0f%%", vol), nil
	//return icon, nil
}
