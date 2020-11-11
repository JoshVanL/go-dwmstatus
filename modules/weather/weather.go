package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
)

const (
	envKey = "GO_DWMSTATUS_WEATHER_API"
)

type Response struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func Weather(h *handler.Handler, s *string) error {
	key := os.Getenv(envKey)
	if len(key) == 0 {
		fmt.Fprintf(os.Stdout, "weather disabled, env %q empty\n", envKey)
		return nil
	}

	update := func() bool {
		temp, err := request(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to request weather: %s\n", err)
			return false
		}

		*s = temp
		h.Tick()
		return true
	}

	ticker := time.NewTicker(time.Second * 5)

	go func() {
		for !update() {
			<-ticker.C
		}

		ticker = time.NewTicker(time.Minute * 10)

		for update() {
			<-ticker.C
		}
	}()

	return nil
}

func request(key string) (string, error) {
	resp, err := http.DefaultClient.Get("http://api.openweathermap.org/data/2.5/weather?q=London&appid=" + key + "&units=metric")
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	respJSON := new(Response)
	if err := json.Unmarshal(b, respJSON); err != nil {
		return "", err
	}

	return fmt.Sprintf(" %.2f°C", respJSON.Main.Temp), nil
}
