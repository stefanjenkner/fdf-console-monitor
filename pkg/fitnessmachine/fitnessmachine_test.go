package fitnessmachine

import (
	"bytes"
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/go-ble/ble"
	"github.com/stefanjenkner/fdf-console-monitor/pkg/events"
)

type NotifierMock struct {
	buffer  *bytes.Buffer
	context *context.Context
}

func TestFitnessMachine_NewSerialMonitor(t *testing.T) {
	want := &FitnessMachine{name: "any name"}
	if got := NewFitnessMachine("any name"); !reflect.DeepEqual(got, want) {
		t.Errorf("NewFitnessMachine() = %v, want %v", got, want)
	}
}

func TestFitnessMachine_rowerDataNotifyHandler(t *testing.T) {
	notifierContext, cancelFunc := context.WithCancel(context.Background())
	var buffer bytes.Buffer
	var notifierMock ble.Notifier = &NotifierMock{
		context: &notifierContext,
		buffer:  &buffer,
	}

	f := &FitnessMachine{}
	go f.rowerDataNotifyHandler(nil, notifierMock)
	time.Sleep(1 * time.Second)
	uint8Ptr := func(v uint8) *uint8 {
		return &v
	}
	uint16Ptr := func(v uint16) *uint16 {
		return &v
	}
	f.OnData(events.DataEvent{
		ElapsedTime:         45,
		Level:               0,
		Distance:            uint16Ptr(123),
		Time500mSplit:       uint16Ptr(115),
		Strokes:             uint16Ptr(23),
		StrokesPerMinute:    uint8Ptr(31),
		WattsPreviousStroke: uint16Ptr(105),
		CaloriesPerHour:     uint16Ptr(987),
	})
	time.Sleep(1 * time.Second)

	(cancelFunc)()

	wanted := []byte{0x2C, 0x08, 62, 23, 0, 123, 0, 0, 115, 0, 105, 0, 45, 0}
	if got := buffer.Bytes(); !bytes.Equal(got, wanted) {
		t.Errorf("wanted: %v, got: %v", wanted, got)
	}
}

func (n *NotifierMock) Context() context.Context {
	return *n.context
}

func (n *NotifierMock) Write(b []byte) (int, error) {
	return (*n.buffer).Write(b)
}

func (n *NotifierMock) Close() error {
	panic("implement me")
}

func (n *NotifierMock) Cap() int {
	panic("implement me")
}
