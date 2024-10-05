package mocks

import (
	"github.com/stefanjenkner/fdf-console-monitor/internal/events"
)

type MockObserver struct {
	DataEvents         []events.DataEvent
	StatusChangeEvents []events.StatusChangeEvent
}

func NewMockObserver() *MockObserver {
	observer := &MockObserver{
		DataEvents:         make([]events.DataEvent, 0),
		StatusChangeEvents: make([]events.StatusChangeEvent, 0),
	}
	return observer
}

func (m *MockObserver) OnData(event events.DataEvent) {
	m.DataEvents = append(m.DataEvents, event)
}

func (m *MockObserver) OnStatusChange(event events.StatusChangeEvent) {
	m.StatusChangeEvents = append(m.StatusChangeEvents, event)
}
