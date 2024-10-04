package observer

import (
	"github.com/stefanjenkner/fdf-console-monitor/pkg/events"
)

type Observer interface {
	OnData(event events.DataEvent)
	OnStatusChange(event events.StatusChangeEvent)
}
