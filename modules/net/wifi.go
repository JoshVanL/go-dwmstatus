package net

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/utils"
)

const (
	wirelessFilePath = "/proc/net/wireless"
)

func Wifi(h *handler.Handler, s *string) error {
	ticker := time.NewTicker(time.Second * 2)

	go func() {
		for {
			w, err := readWireless()
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get wireless info: %s", err)
			}

			if len(w) > 0 {
				*s = fmt.Sprintf(" %s%% ", w)

				////if n >= 70 {
				////	block.Color = "#0050b8"
				////} else if n >= 50 {
				//if n >= 50 {
				//	//block.Color = "#aaddaa"
				//} else if n >= 30 {
				//	//block.Color = "#a65f3d"
				//	//block.Color = "#ffae88"
				//} else {
				//	//block.Color = "#ff0000"
				//}

			} else {
				*s = ""
			}
			h.Tick()

			<-ticker.C
		}
	}()

	return nil
}

func readWireless() (string, error) {
	b, err := utils.ReadFile(wirelessFilePath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(b), "\n")
	if len(lines) < 3 {
		return "", nil
	}

	fields := strings.Fields(lines[2])
	return strings.TrimSuffix(fields[2], "."), nil
}
