package watcher

import (
	"sync"
	"syscall"
	"unsafe"
)

const (
	sys_IN_MODIFY uint32 = syscall.IN_MODIFY

	SIGRT_MIN = syscall.Signal(34)
	SIGRT_MAX = syscall.Signal(65)
)

var (
	RealTimeSignals = map[string]syscall.Signal{
		"RTMIN":    SIGRT_MIN,
		"RTMIN+1":  SIGRT_MIN + 1,
		"RTMIN+2":  SIGRT_MIN + 2,
		"RTMIN+3":  SIGRT_MIN + 3,
		"RTMIN+4":  SIGRT_MIN + 4,
		"RTMIN+5":  SIGRT_MIN + 5,
		"RTMIN+6":  SIGRT_MIN + 6,
		"RTMIN+7":  SIGRT_MIN + 7,
		"RTMIN+8":  SIGRT_MIN + 8,
		"RTMIN+9":  SIGRT_MIN + 9,
		"RTMIN+10": SIGRT_MIN + 10,
		"RTMIN+11": SIGRT_MIN + 11,
		"RTMIN+12": SIGRT_MIN + 12,
		"RTMIN+13": SIGRT_MIN + 13,
		"RTMIN+14": SIGRT_MIN + 14,
		"RTMIN+15": SIGRT_MIN + 15,
		"RTMAX-14": SIGRT_MAX - 14,
		"RTMAX-13": SIGRT_MAX - 13,
		"RTMAX-12": SIGRT_MAX - 12,
		"RTMAX-11": SIGRT_MAX - 11,
		"RTMAX-10": SIGRT_MAX - 10,
		"RTMAX-9":  SIGRT_MAX - 9,
		"RTMAX-8":  SIGRT_MAX - 8,
		"RTMAX-7":  SIGRT_MAX - 7,
		"RTMAX-6":  SIGRT_MAX - 6,
		"RTMAX-5":  SIGRT_MAX - 5,
		"RTMAX-4":  SIGRT_MAX - 4,
		"RTMAX-3":  SIGRT_MAX - 3,
		"RTMAX-2":  SIGRT_MAX - 2,
		"RTMAX-1":  SIGRT_MAX - 1,
		"RTMAX":    SIGRT_MAX,
	}
)

type Watcher struct {
	fd       int
	watching map[int32]chan struct{}
	mu       sync.Mutex
}

func New() (*Watcher, error) {
	fd, err := syscall.InotifyInit()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		fd:       fd,
		watching: make(map[int32]chan struct{}),
	}

	go w.run()

	return w, nil
}

func (w *Watcher) run() {
	var buf [syscall.SizeofInotifyEvent * 1024]byte
	for {
		_, err := syscall.Read(w.fd, buf[:])
		if err != nil {
			continue
		}

		raw := (*syscall.InotifyEvent)(unsafe.Pointer(&buf))

		if (raw.Mask & sys_IN_MODIFY) != sys_IN_MODIFY {
			continue
		}

		w.mu.Lock()
		ch, ok := w.watching[raw.Wd]
		w.mu.Unlock()

		if !ok {
			continue
		}

		ch <- struct{}{}
	}
}

func (w *Watcher) Add(path string) (<-chan struct{}, error) {
	ch := make(chan struct{})

	wd, err := syscall.InotifyAddWatch(w.fd, path, sys_IN_MODIFY)
	if err != nil {
		return nil, err
	}

	w.mu.Lock()
	w.watching[int32(wd)] = ch
	w.mu.Unlock()

	return ch, nil
}
