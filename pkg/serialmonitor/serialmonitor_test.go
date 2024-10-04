package serialmonitor

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stefanjenkner/fdf-console-monitor/mocks"
	"github.com/stefanjenkner/fdf-console-monitor/pkg/events"
	"github.com/stefanjenkner/fdf-console-monitor/pkg/observer"
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
		*events.NewDataEventBuilder(4, 4).SetDistance(7).SetTime500mSplit(268).SetStrokes(1).SetStrokesPerMinute(14).SetWattsPreviousStroke(108).SetCaloriesPerHour(670).Build(),
		*events.NewDataEventBuilder(6, 4).SetDistance(14).SetTime500mSplit(163).SetStrokes(2).SetStrokesPerMinute(28).SetWattsPreviousStroke(105).SetCaloriesPerHour(659).Build(),
		*events.NewDataEventBuilder(8, 4).SetDistance(21).SetTime500mSplit(148).SetStrokes(3).SetStrokesPerMinute(29).SetWattsPreviousStroke(109).SetCaloriesPerHour(674).Build(),
		*events.NewDataEventBuilder(2, 4).SetRemainingDistance(0).SetTime500mAverage(0).SetWattsAverage(0).SetCaloriesTotal(0).Build(),
		*events.NewDataEventBuilder(5, 4).SetDistance(8).SetTime500mSplit(319).SetStrokes(1).SetStrokesPerMinute(11).SetWattsPreviousStroke(106).SetCaloriesPerHour(663).Build(),
		*events.NewDataEventBuilder(1810, 4).SetRemainingDistance(6015).SetTime500mAverage(153).SetWattsAverage(109).SetCaloriesTotal(400).Build(),
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
		{StatusChange: events.Started},
		{StatusChange: events.Reset},
		{StatusChange: events.PausedOrStopped},
		{StatusChange: events.Resumed},
		{StatusChange: events.PausedOrStopped},
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
