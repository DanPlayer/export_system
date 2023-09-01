package timewheel

import (
	"github.com/rfyiamcool/go-timewheel"
	"time"
)

var Client = NewClient()

func NewClient() *timewheel.TimeWheel {
	tw, err := timewheel.NewTimeWheel(1*time.Second, 360, timewheel.TickSafeMode())
	if err != nil {
		panic(err)
	}
	return tw
}
