package handler

// #cgo LDFLAGS: -lX11 -lasound
// #include <X11/Xlib.h>
import "C"

import (
	"fmt"
	"sync"

	"github.com/joshvanl/go-dwmstatus/errors"
	"github.com/joshvanl/go-dwmstatus/sysinfo"
	"github.com/joshvanl/go-dwmstatus/watcher"
)

var (
	dpy = C.XOpenDisplay(nil)
)

type Handler struct {
	watcher *watcher.Watcher
	sysinfo *sysinfo.SysInfo

	mu      sync.Mutex
	strings []*string
}

func New() (*Handler, error) {
	w, err := watcher.New()
	if err != nil {
		return nil, err
	}

	s, err := sysinfo.New()
	if err != nil {
		return nil, err
	}

	h := &Handler{
		watcher: w,
		sysinfo: s,
	}

	if dpy == nil {
		h.Must(fmt.Errorf("failed to open display: %v", dpy))
	}

	go h.signalHandler()

	return h, nil
}

func (h *Handler) NewModule() *string {
	h.mu.Lock()
	defer h.mu.Unlock()

	s := new(string)
	h.strings = append(h.strings, s)

	return s
}

func (h *Handler) Tick() {
	h.mu.Lock()
	defer h.mu.Unlock()

	var output string
	for _, s := range h.strings {
		output += *s
	}

	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), C.CString(output))
	C.XSync(dpy, 1)
}

func (h *Handler) SysInfo() *sysinfo.SysInfo {
	return h.sysinfo
}

func (h *Handler) Must(err error) {
	if err == nil {
		return
	}

	errors.Kill(fmt.Errorf("go-dwmstatus was killed: %v\n", err))
}

func (h *Handler) signalHandler() {
	ch := errors.NewSignalHandler()
	<-ch
	h.Must(fmt.Errorf("got signal interupt"))
}

func (h *Handler) Watcher() *watcher.Watcher {
	return h.watcher
}
