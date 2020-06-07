package temp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
	"github.com/joshvanl/go-dwmstatus/modules/utils"
)

const (
	thermalDir  = "/sys/class/thermal"
	thermalType = "x86_pkg_temp"
)

func Temp(h *handler.Handler, s *string) error {
	files, err := ioutil.ReadDir(thermalDir)
	if err != nil {
		return err
	}

	var path string
	for _, f := range files {
		b, err := utils.ReadFile(filepath.Join(thermalDir, f.Name(), "type"))
		if err != nil {
			return err
		}

		if string(b) == thermalType {
			path = filepath.Join(thermalDir, f.Name(), "temp")
			break
		}
	}

	if path == "" {
		return fmt.Errorf("failed to find thermal with type %s",
			thermalType)
	}

	//normalColor := block.Color

	ticker := time.NewTicker(time.Second * 5)

	go func() {
		for {
			b, err := utils.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read temp file: %s\n", err)
			}

			temp, err := strconv.ParseFloat(string(b), 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to get temp: %s\n", err)
			}

			temp = temp / 1000

			if temp > 90 {
				//block.Color = "#ff3333"
			} else if temp > 70 {
				//block.Color = "#ffaa33"
			} else {
				//block.Color = normalColor
			}

			*s = fmt.Sprintf("%1.f°C", temp)
			h.Tick()

			<-ticker.C
		}
	}()

	return nil
}
