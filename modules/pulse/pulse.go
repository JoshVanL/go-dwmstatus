package pulse

import (
	"fmt"

	"github.com/auroralaboratories/pulse"
	"mrogalski.eu/go/pulseaudio"
)

type pulseHandler struct {
	events <-chan struct{}
	*pulse.Client
}

func getHandler() (*pulseHandler, error) {
	client, err := pulse.NewClient("go-dwnstatus-client")
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate shared pulse client: %s", err)
	}

	subscriptionClient, err := pulseaudio.NewClient()
	if err != nil {
		return nil, err
	}

	ch, err := subscriptionClient.Updates()
	if err != nil {
		return nil, err
	}

	sharedHandler := &pulseHandler{
		events: ch,
		Client: client,
	}

	return sharedHandler, nil
}
