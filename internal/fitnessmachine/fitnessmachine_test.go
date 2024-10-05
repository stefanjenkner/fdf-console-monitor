package fitnessmachine

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-ble/ble"
	"github.com/stefanjenkner/fdf-console-monitor/internal/events"
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
	f.OnData(*events.NewDataEvent(45, 0,
		events.WithDistance(123),
		events.WithTime500mSplit(115),
		events.WithStrokes(23),
		events.WithStrokesPerMinute(31),
		events.WithWattsPreviousStroke(105),
		events.WithCaloriesPerHour(987),
	))
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
	write, err := (*n.buffer).Write(b)
	if err != nil {
		return 0, fmt.Errorf("write: unexpected error: %w", err)
	}
	return write, nil
}

func (n *NotifierMock) Close() error {
	panic("implement me")
}

func (n *NotifierMock) Cap() int {
	panic("implement me")
}
