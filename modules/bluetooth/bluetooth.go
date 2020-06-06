package bluetooth

import (
	"fmt"
	"os"
	"reflect"

	"github.com/joshvanl/go-dwmstatus/handler"

	//"github.com/godbus/dbus/introspect"
	//"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
)

const (
	iface = "hci0"
)

var (
	watchDevices = []string{
		"9A:BC:04:01:DD:DA", // E18 Plus
	}
)

type bluetoothHandler struct {
	powerOn          bool
	connectedDevices map[string]struct{}

	s *string

	adapter *adapter.Adapter1
	handler *handler.Handler

	selectCaseWatchers []reflect.SelectCase
	names              []string
}

func Bluetooth(h *handler.Handler, s *string) error {
	a, err := adapter.GetAdapter(iface)
	if err != nil {
		return err
	}

	bluetoothHandler := &bluetoothHandler{
		powerOn:          a.Properties.Powered,
		connectedDevices: make(map[string]struct{}),
		s:                s,
		adapter:          a,
		handler:          h,
	}

	if err := bluetoothHandler.watch(); err != nil {
		return err
	}

	return nil
}

func (b *bluetoothHandler) watch() error {
	adapterCh, err := b.adapter.WatchProperties()
	if err != nil {
		return err
	}

	b.names = make([]string, 0)
	b.selectCaseWatchers = make([]reflect.SelectCase, 0)

	b.selectCaseWatchers = append(b.selectCaseWatchers, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(adapterCh)})
	b.names = append(b.names, b.adapter.Properties.Name)

	fmt.Printf("watching bluetooth adapter: %s\n", b.adapter.Properties.Name)

	devices, err := b.adapter.GetDevices()
	if err != nil {
		return err
	}

	for _, d := range devices {
		for _, address := range watchDevices {
			if d.Properties.Address == address || d.Properties.Connected {
				ch, err := d.WatchProperties()
				if err != nil {
					return err
				}

				fmt.Printf("watching bluetooth device: %s\n", d.Properties.Name)

				b.selectCaseWatchers = append(b.selectCaseWatchers, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
				b.names = append(b.names, d.Properties.Name)
				if d.Properties.Connected {
					b.connectedDevices[d.Properties.Name] = struct{}{}
				}
			}
		}
	}

	go func() {
		for {
			if err := b.nextEvent(); err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
			}
		}
	}()

	return nil
}

func (b *bluetoothHandler) nextEvent() error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "recovered bluetooth panic: %s\n", r)
		}
	}()

	chosen, value, _ := reflect.Select(b.selectCaseWatchers)

	prop, ok := value.Interface().(*bluez.PropertyChanged)
	if !ok {
		return fmt.Errorf("failed to get propery value: %s\n", value.String())
	}

	switch prop.Interface {
	case adapter.Adapter1Interface:
		if prop.Name == "Powered" {
			b.powerOn = prop.Value.(bool)
		}

	case device.Device1Interface:
		if prop.Name != "Connected" {
			return nil
		}

		if prop.Value.(bool) {
			b.connectedDevices[b.names[chosen]] = struct{}{}
		} else {
			delete(b.connectedDevices, b.names[chosen])
		}

	default:
		return fmt.Errorf("unrecognised property interface: %s\n", prop.Interface)
	}

	b.update()
	b.handler.Tick()

	return nil
}

func (b *bluetoothHandler) update() {
	var output string
	for name := range b.connectedDevices {
		output += name + " "
	}

	if b.powerOn {
		output += " "
	}

	*b.s = output
}
