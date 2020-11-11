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

	go func() {
		now := time.Now()
		time.Sleep(
			time.Until(
				time.Date(
					now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, loc).
					Add(time.Minute)))

		ticker := time.NewTicker(time.Minute)

		for {
			*s = getTimeString(time.Now().In(loc))
			h.Tick()

			<-ticker.C
		}
	}()

	return nil
}

func getTimeString(t time.Time) string {
	return fmt.Sprintf(`%s %d %s %d %02d:%02d`,
		t.Format("Mon"), t.Day(), t.Month().String()[:3], t.Year(), t.Hour(), t.Minute())
}
