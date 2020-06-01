package net

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/mickep76/netlink"
)

const (
	netDir = "/sys/class/net"
)

func IFace(h *handler.Handler, s *string) error {
	ifaceHandler, err := getNetHandler()
	if err != nil {
		return err
	}

	go func() {
		for {
			ips, ok, err := getIPs(ifaceHandler.ifaces)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get iface addresses: %s", err)
				time.Sleep(time.Second)
				continue
			}

			// If iface is running but address not ready, try again in half a second
			if !ok {
				time.Sleep(time.Second / 2)
				continue
			}

			*s = ""
			for name, ip := range ips {
				*s += name + ":" + ip
			}

			h.Tick()
			ifaceHandler.cond.Wait()
		}
	}()

	return nil
}

// getIPs will return a map of interfaces to their IP addresses. Returns false
// if interface is running, but has not been assigned IP
func getIPs(ifaces []netlink.Interface) (map[string]string, bool, error) {
	ips := make(map[string]string)
	for _, iface := range ifaces {
		// ignore loopback device
		if iface.Name == "lo" {
			continue
		}

		addrs, err := iface.NetInterface.Addrs()
		if err != nil {
			return nil, false, fmt.Errorf("failed to get iface %q addresses: %s",
				iface.Name, err)
		}

		var ip net.IP

	LOOP:
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
				break LOOP
			case *net.IPAddr:
				ip = v.IP
				break LOOP
			}
		}

		if ip == nil {
			continue
		}

		if ip.To4() != nil {
			ips[iface.Name] = ip.String()
			continue
		}

		if iface.Flags&netlink.FlagRunning != 0 {
			return nil, false, nil
		}
	}

	return ips, true, nil
}
