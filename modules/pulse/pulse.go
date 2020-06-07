package pulse

import (
	"fmt"
	"sync"

	"github.com/auroralaboratories/pulse"
)

var (
	sharedHandler *pulseHandler
)

type pulseHandler struct {
	mu     sync.Mutex
	events <-chan pulse.EventType

	client *pulse.Client
}

func getHandler() (*pulseHandler, error) {
	if sharedHandler != nil {
		return sharedHandler, nil
	}

	client, err := pulse.NewClient("go-dwnstatus-client")
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate shared pulse client: %s", err)
	}

	subscriptionClient, err := pulse.NewClient("go-dwnstatus-watcher")
	if err != nil {
		return nil, err
	}

	ch := subscriptionClient.Subscribe(pulse.AllEvent)

	sharedHandler = &pulseHandler{
		events: ch,
		client: client,
	}

	return sharedHandler, nil
}

func (p *pulseHandler) GetSinks() ([]pulse.Sink, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.client.GetSinks()
}

func (p *pulseHandler) GetSources() ([]pulse.Source, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.client.GetSources()
}
