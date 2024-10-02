package fitnessmachine

import (
	"bytes"
	"context"
	"fdf-console-monitor/internal/events"
	"reflect"
	"testing"
	"time"

	"github.com/go-ble/ble"
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
	uint64Ptr := func(v uint64) *uint64 {
		return &v
	}
	f.OnData(events.DataEvent{
		ElapsedTime:         45,
		Level:               0,
		Distance:            uint64Ptr(123),
		Time500mSplit:       uint64Ptr(115),
		Strokes:             uint64Ptr(23),
		StrokesPerMinute:    uint64Ptr(31),
		WattsPreviousStroke: uint64Ptr(105),
		CaloriesPerHour:     uint64Ptr(987),
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
