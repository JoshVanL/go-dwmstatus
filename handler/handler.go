package handler

// #cgo LDFLAGS: -lX11 -lasound
// #include <X11/Xlib.h>
import "C"

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
		return nil, fmt.Errorf("failed to open display: %v", dpy)
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

func (h *Handler) signalHandler() {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)

	recSig := <-sig

	Kill(fmt.Errorf("go-dwmstatus was killed: got signal: %s", recSig))
}

func (h *Handler) Watcher() *watcher.Watcher {
	return h.watcher
}

func Kill(err error) {
	defer func() {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(0)
	}()

	if err := os.MkdirAll("/home/josh/.cache/go-dwmstatus", 0755); err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	f, ferr := os.OpenFile("/home/josh/.cache/go-dwmstatus/err.log", os.O_CREATE|os.O_WRONLY, 0666)
	if ferr != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}

	fmt.Fprint(f, err)
	f.Close()
}
