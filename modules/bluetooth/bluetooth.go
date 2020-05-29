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
	var selectCaseWatchers []reflect.SelectCase
	var names []string

	adapterCh, err := b.adapter.WatchProperties()
	if err != nil {
		return err
	}

	selectCaseWatchers = append(selectCaseWatchers, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(adapterCh)})
	names = append(names, b.adapter.Properties.Name)

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

				selectCaseWatchers = append(selectCaseWatchers, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)})
				names = append(names, d.Properties.Name)
				if d.Properties.Connected {
					b.connectedDevices[d.Properties.Name] = struct{}{}
				}
			}
		}
	}

	go func() {
		for {
			chosen, value, _ := reflect.Select(selectCaseWatchers)

			prop, ok := value.Interface().(*bluez.PropertyChanged)
			if !ok {
				fmt.Fprintf(os.Stderr, "failed to get propery value: %s\n", value.String())
				continue
			}
			fmt.Printf("%#+v\n", prop)

			switch prop.Interface {
			case adapter.Adapter1Interface:
				if prop.Name == "Powered" {
					b.powerOn = prop.Value.(bool)
				}

			case device.Device1Interface:
				if prop.Name != "Connected" {
					continue
				}

				if prop.Value.(bool) {
					b.connectedDevices[names[chosen]] = struct{}{}
				} else {
					delete(b.connectedDevices, names[chosen])
				}

			default:
				fmt.Fprintf(os.Stderr, "unrecognised property interface: %s\n", prop.Interface)
				continue
			}

			b.update()
			b.handler.Tick()
		}
	}()

	return nil
}

func (b *bluetoothHandler) update() {
	var output string
	for name := range b.connectedDevices {
		output += name
	}

	if b.powerOn {
		output += " "
	}

	*b.s = output
}
