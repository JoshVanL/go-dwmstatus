package net

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/mickep76/netlink"
)

const (
	sysNetPath = "/sys/class/net"
)

var (
	sharedNetHandler *netHandler
)

type netHandler struct {
	cond *sync.Cond

	ifaces []netlink.Interface
}

func getNetHandler() (*netHandler, error) {
	if sharedNetHandler != nil {
		return sharedNetHandler, nil
	}

	mu := new(sync.Mutex)

	sharedNetHandler = &netHandler{
		cond: sync.NewCond(mu),
	}
	sharedNetHandler.cond.L.Lock()

	if err := sharedNetHandler.updateIFaces(); err != nil {
		return nil, err
	}

	clnt, err := netlink.Dial(netlink.NetlinkRoute, netlink.RtmGrpLink)
	if err != nil {
		return nil, err
	}

	if err := clnt.Bind(); err != nil {
		return nil, err
	}

	// TODO
	// defer clnt.Close()

	go func() {
		for {
			_, err := clnt.Receive()
			if err != nil {
				log.Fatal(err)
			}

			if err := sharedNetHandler.updateIFaces(); err != nil {
				fmt.Fprintf(os.Stderr, "failed to get interfaces: %s\n", err)
			}

			sharedNetHandler.cond.Broadcast()
		}
	}()

	return sharedNetHandler, nil
}

func (n *netHandler) updateIFaces() error {
	ifs, err := netlink.Interfaces()
	if err != nil {
		return err
	}

	n.ifaces = ifs

	return nil
}
