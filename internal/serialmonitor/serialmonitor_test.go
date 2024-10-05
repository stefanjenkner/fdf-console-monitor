package serialmonitor

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stefanjenkner/fdf-console-monitor/internal/events"
	"github.com/stefanjenkner/fdf-console-monitor/internal/mocks"
	"github.com/stefanjenkner/fdf-console-monitor/internal/observer"
)

func TestSerialMonitor_RunCallsObserverForDataEvents(t *testing.T) {
	bufferString := bytes.NewBufferString("")
	// session one
	bufferString.WriteString("A8000040000710428014108067004\r\n")
	bufferString.WriteString("A8000060001410243028105065904\r\n")
	bufferString.WriteString("A8000080002110228029109067404\r\n")
	// reset
	bufferString.WriteString("R\r\n")
	// session two
	bufferString.WriteString("A8000020000010000000000000004\r\n")
	bufferString.WriteString("A8000050000810519011106066304\r\n")
	// ... skip some
	bufferString.WriteString("A8030100601510233000109040004\r\n")
	// connection check
	bufferString.WriteString("W\r\n")

	mockSerialPort, port := mocks.NewMockSerialPort(bufferString)
	serialMonitor := SerialMonitor{
		portName:  "/dev/mocked/serial/port",
		port:      &port,
		observers: map[observer.Observer]struct{}{},
	}
	mockObserver := mocks.NewMockObserver()
	serialMonitor.AddObserver(mockObserver)
	serialMonitor.Run()

	wantedDataEvents := []events.DataEvent{
		*events.NewDataEvent(4, 4, events.WithDistance(7), events.WithTime500mSplit(268),
			events.WithStrokes(1), events.WithStrokesPerMinute(14), events.WithWattsPreviousStroke(108),
			events.WithCaloriesPerHour(670),
		),
		*events.NewDataEvent(6, 4, events.WithDistance(14), events.WithTime500mSplit(163),
			events.WithStrokes(2), events.WithStrokesPerMinute(28), events.WithWattsPreviousStroke(105),
			events.WithCaloriesPerHour(659),
		),
		*events.NewDataEvent(8, 4, events.WithDistance(21), events.WithTime500mSplit(148),
			events.WithStrokes(3), events.WithStrokesPerMinute(29), events.WithWattsPreviousStroke(109),
			events.WithCaloriesPerHour(674),
		),
		*events.NewDataEvent(2, 4, events.WithRemainingDistance(0), events.WithTime500mAverage(0),
			events.WithWattsAverage(0), events.WithCaloriesTotal(0),
		),
		*events.NewDataEvent(5, 4, events.WithDistance(8), events.WithTime500mSplit(319),
			events.WithStrokes(1), events.WithStrokesPerMinute(11), events.WithWattsPreviousStroke(106),
			events.WithCaloriesPerHour(663),
		),
		*events.NewDataEvent(1810, 4, events.WithRemainingDistance(6015), events.WithTime500mAverage(153),
			events.WithWattsAverage(109),
			events.WithCaloriesTotal(400),
		),
	}

	if got := len(mockObserver.DataEvents); len(wantedDataEvents) != got {
		t.Errorf("# dataEvents = %v, wantedDataEvents %v", got, len(wantedDataEvents))
		t.FailNow()
	}

	for i := 0; i < len(mockObserver.DataEvents); i++ {
		if got := mockObserver.DataEvents[i]; !reflect.DeepEqual(got, wantedDataEvents[i]) {
			t.Errorf("dataEvents() = %+v, wantedDataEvents %+v", got, wantedDataEvents[i])
		}
	}

	if got := mockSerialPort.Closed; got != 1 {
		t.Errorf("mockSerialPort.closed = %v, wantedDataEvents 1", got)
	}
}

func TestSerialMonitor_RunCallsObserverForStatusChangeEvents(t *testing.T) {
	bufferString := bytes.NewBufferString("")
	// session one
	bufferString.WriteString("A8000040000710428014108067004\r\n")
	bufferString.WriteString("A8000060001410243028105065904\r\n")
	bufferString.WriteString("A8000080002110228029109067404\r\n")
	// reset
	bufferString.WriteString("R\r\n")
	// session two
	bufferString.WriteString("A8000020000010000000000000004\r\n")
	bufferString.WriteString("A8000050000810519011106066304\r\n")
	// ... skip some
	bufferString.WriteString("A8030100601510233000109040004\r\n")
	// connection check
	bufferString.WriteString("W\r\n")

	mockSerialPort, port := mocks.NewMockSerialPort(bufferString)
	serialMonitor := SerialMonitor{
		portName:  "/dev/mocked/serial/port",
		port:      &port,
		observers: map[observer.Observer]struct{}{},
	}
	mockObserver := mocks.NewMockObserver()
	serialMonitor.AddObserver(mockObserver)
	serialMonitor.Run()

	wantedStatusChangeEvents := []events.StatusChangeEvent{
		*events.NewStatusChangeEvent(events.Started),
		*events.NewStatusChangeEvent(events.Reset),
		*events.NewStatusChangeEvent(events.PausedOrStopped),
		*events.NewStatusChangeEvent(events.Resumed),
		*events.NewStatusChangeEvent(events.PausedOrStopped),
	}
	if got := len(mockObserver.StatusChangeEvents); len(wantedStatusChangeEvents) != got {
		t.Errorf("# statusChangeEvents = %v, wantedDataEvents %v", got, len(wantedStatusChangeEvents))
		t.FailNow()
	}
	for i := 0; i < len(mockObserver.StatusChangeEvents); i++ {
		if got := mockObserver.StatusChangeEvents[i]; !reflect.DeepEqual(got, wantedStatusChangeEvents[i]) {
			t.Errorf("statusChangeEvents() = %+v, wantedDataEvents %+v", got, wantedStatusChangeEvents[i])
		}
	}
	if got := mockSerialPort.Closed; got != 1 {
		t.Errorf("mockSerialPort.closed = %v, wantedDataEvents 1", got)
	}
}

func TestSerialMonitor_NewSerialMonitor(t *testing.T) {
	want := &SerialMonitor{
		portName:  "/any/port/name",
		observers: make(map[observer.Observer]struct{}),
	}
	if got := NewSerialMonitor("/any/port/name"); !reflect.DeepEqual(got, want) {
		t.Errorf("NewSerialMonitor() = %v, want %v", got, want)
	}
}
