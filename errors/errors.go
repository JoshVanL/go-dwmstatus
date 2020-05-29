package errors

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func NewSignalHandler() <-chan int {
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGKILL, syscall.SIGINT)

	ch := make(chan int)

	go func() {
		for s := range sig {
			switch s {
			case syscall.SIGKILL, syscall.SIGINT:
				ch <- -1

			default:
				break
			}
		}
	}()

	return ch
}

func Kill(err error) {
	defer func() {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
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
