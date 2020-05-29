package datetime

import (
	"fmt"
	"time"

	"github.com/joshvanl/go-dwmstatus/handler"
)

func DateTime(h *handler.Handler, s *string) error {
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return err
	}

	*s = getTimeString(time.Now().In(loc))
	h.Tick()

	now := time.Now()
	time.Sleep(now.Truncate(time.Second).Add(time.Second).Sub(now))
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			*s = getTimeString(time.Now().In(loc))
			h.Tick()

			<-ticker.C
		}
	}()

	return nil
}

func getTimeString(t time.Time) string {
	return fmt.Sprintf("%s %d %s %d %s",
		t.Format("Mon"), t.Day(), t.Month().String()[:3], t.Year(), t.Format("15:04:05"))
}
