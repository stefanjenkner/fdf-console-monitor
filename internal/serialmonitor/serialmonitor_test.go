package serialmonitor

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"fdf-console-monitor/internal/events"

	"go.bug.st/serial"
)

type MockSerialPort struct {
	readBuffer  bytes.Buffer
	writeBuffer bytes.Buffer
	closed      int
}

type MockObserver struct {
	dataEvents         []events.DataEvent
	statusChangeEvents []events.StatusChangeEvent
}

func TestSerialMonitor_Run(t *testing.T) {
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

	mockSerialPort := MockSerialPort{
		readBuffer:  *bufferString,
		writeBuffer: bytes.Buffer{},
	}
	var port serial.Port = &mockSerialPort
	serialMonitor := SerialMonitor{
		portName:  "/dev/mocked/serial/port",
		port:      &port,
		observers: map[events.Observer]struct{}{},
	}
	observer := &MockObserver{
		dataEvents:         make([]events.DataEvent, 0),
		statusChangeEvents: make([]events.StatusChangeEvent, 0),
	}
	serialMonitor.AddObserver(observer)
	serialMonitor.Run()

	uint8Ptr := func(v uint8) *uint8 {
		return &v
	}
	uint16Ptr := func(v uint16) *uint16 {
		return &v
	}
	wantedDataEvents := []events.DataEvent{
		{
			ElapsedTime:         4,
			Level:               4,
			Distance:            uint16Ptr(7),
			Time500mSplit:       uint16Ptr(268),
			Strokes:             uint16Ptr(1),
			StrokesPerMinute:    uint8Ptr(14),
			WattsPreviousStroke: uint16Ptr(108),
			CaloriesPerHour:     uint16Ptr(670),
		}, {
			ElapsedTime:         6,
			Level:               4,
			Distance:            uint16Ptr(14),
			Time500mSplit:       uint16Ptr(163),
			Strokes:             uint16Ptr(2),
			StrokesPerMinute:    uint8Ptr(28),
			WattsPreviousStroke: uint16Ptr(105),
			CaloriesPerHour:     uint16Ptr(659),
		}, {
			ElapsedTime:         8,
			Level:               4,
			Distance:            uint16Ptr(21),
			Time500mSplit:       uint16Ptr(148),
			Strokes:             uint16Ptr(3),
			StrokesPerMinute:    uint8Ptr(29),
			WattsPreviousStroke: uint16Ptr(109),
			CaloriesPerHour:     uint16Ptr(674),
		}, {
			ElapsedTime:       2,
			Level:             4,
			RemainingDistance: uint16Ptr(0),
			Time500mAverage:   uint16Ptr(0),
			WattsAverage:      uint16Ptr(0),
			CaloriesTotal:     uint16Ptr(0),
		}, {
			ElapsedTime:         5,
			Level:               4,
			Distance:            uint16Ptr(8),
			Time500mSplit:       uint16Ptr(319),
			Strokes:             uint16Ptr(1),
			StrokesPerMinute:    uint8Ptr(11),
			WattsPreviousStroke: uint16Ptr(106),
			CaloriesPerHour:     uint16Ptr(663),
		}, {
			ElapsedTime:       1810,
			Level:             4,
			RemainingDistance: uint16Ptr(6015),
			Time500mAverage:   uint16Ptr(153),
			WattsAverage:      uint16Ptr(109),
			CaloriesTotal:     uint16Ptr(400),
		},
	}

	if got := len(observer.dataEvents); len(wantedDataEvents) != got {
		t.Errorf("# dataEvents = %v, wantedDataEvents %v", got, len(wantedDataEvents))
		t.FailNow()
	}

	for i := 0; i < len(observer.dataEvents); i++ {
		if got := observer.dataEvents[i]; !reflect.DeepEqual(got, wantedDataEvents[i]) {
			t.Errorf("dataEvents() = %+v, wantedDataEvents %+v", got, wantedDataEvents[i])
		}
	}

	wantedStatusChangeEvents := []events.StatusChangeEvent{
		{StatusChange: events.Started},
		{StatusChange: events.Reset},
		{StatusChange: events.PausedOrStopped},
		{StatusChange: events.Resumed},
		{StatusChange: events.PausedOrStopped},
	}
	if got := len(observer.statusChangeEvents); len(wantedStatusChangeEvents) != got {
		t.Errorf("# statusChangeEvents = %v, wantedDataEvents %v", got, len(wantedStatusChangeEvents))
		t.FailNow()
	}
	for i := 0; i < len(observer.statusChangeEvents); i++ {
		if got := observer.statusChangeEvents[i]; !reflect.DeepEqual(got, wantedStatusChangeEvents[i]) {
			t.Errorf("statusChangeEvents() = %+v, wantedDataEvents %+v", got, wantedStatusChangeEvents[i])
		}
	}

	if got := mockSerialPort.closed; got != 1 {
		t.Errorf("mockSerialPort.closed = %v, wantedDataEvents 1", got)
	}
}

func TestSerialMonitor_NewSerialMonitor(t *testing.T) {
	want := &SerialMonitor{
		portName:  "/any/port/name",
		observers: make(map[events.Observer]struct{}),
	}
	if got := NewSerialMonitor("/any/port/name"); !reflect.DeepEqual(got, want) {
		t.Errorf("NewSerialMonitor() = %v, want %v", got, want)
	}
}

func (m *MockObserver) OnData(event events.DataEvent) {
	m.dataEvents = append(m.dataEvents, event)
}

func (m *MockObserver) OnStatusChange(event events.StatusChangeEvent) {
	m.statusChangeEvents = append(m.statusChangeEvents, event)
}

func (m *MockSerialPort) Read(p []byte) (n int, err error) {
	return m.readBuffer.Read(p)
}

func (m *MockSerialPort) Write(p []byte) (n int, err error) {
	return m.writeBuffer.Write(p)
}

func (m *MockSerialPort) Close() error {
	m.closed++
	return nil
}

func (m *MockSerialPort) SetMode(_ *serial.Mode) error {
	panic("implement me")
}

func (m *MockSerialPort) Drain() error {
	panic("implement me")
}

func (m *MockSerialPort) ResetInputBuffer() error {
	panic("implement me")
}

func (m *MockSerialPort) ResetOutputBuffer() error {
	panic("implement me")
}

func (m *MockSerialPort) SetDTR(_ bool) error {
	panic("implement me")
}

func (m *MockSerialPort) SetRTS(_ bool) error {
	panic("implement me")
}

func (m *MockSerialPort) GetModemStatusBits() (*serial.ModemStatusBits, error) {
	panic("implement me")
}

func (m *MockSerialPort) SetReadTimeout(_ time.Duration) error {
	panic("implement me")
}

func (m *MockSerialPort) Break(_ time.Duration) error {
	panic("implement me")
}
