package fitnessmachine

import (
	"context"
	"encoding/binary"
	"log"
	"time"

	"fdf-console-monitor/internal/events"
	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

type FitnessMachine struct {
	name            string
	service         *ble.Service
	cancelFunc      *context.CancelFunc
	rowerDataEvents *chan events.DataEvent
}

func NewFitnessMachine(name string) *FitnessMachine {
	return &FitnessMachine{name: name}
}

func (f *FitnessMachine) Start() {
	log.Printf("Starting FitnessMachine: %s\n", f.name)

	device, err := dev.NewDevice("default")
	if err != nil {
		log.Printf("Error creating new device: %+v\n", err)
		return
	}
	ble.SetDefaultDevice(device)
	f.service = ble.NewService(ble.UUID16(0x1826))

	// rower feature
	rowerFeatureChar := ble.NewCharacteristic(ble.UUID16(0x2ACC))
	rowerFeatureChar.HandleRead(ble.ReadHandlerFunc(f.rowerFeatureReadHandler))
	f.service.AddCharacteristic(rowerFeatureChar)

	// rower data
	rowerDataChar := ble.NewCharacteristic(ble.UUID16(0x2AD1))
	rowerDataChar.HandleNotify(ble.NotifyHandlerFunc(f.rowerDataNotifyHandler))
	f.service.AddCharacteristic(rowerDataChar)

	if err = ble.AddService(f.service); err != nil {
		log.Printf("Error adding service: %+v\n", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := ble.AdvertiseNameAndServices(ctx, f.name, f.service.UUID); err != nil {
			log.Printf("Error advertising device and services: %+v\n", err)
		}
	}()

	f.cancelFunc = &cancel
}

func (f *FitnessMachine) Stop() {
	log.Println("Stopping FitnessMachine...")
	if f.cancelFunc != nil {
		(*f.cancelFunc)()
	}
	if err := ble.Stop(); err != nil {
		log.Printf("Error stopping FitnessMachine: %+v\n", err)
		return
	}
	log.Println("Stopped FitnessMachine")
}

func (f *FitnessMachine) rowerDataNotifyHandler(_ ble.Request, n ble.Notifier) {
	log.Println("Subscription started")

	dataEvents := make(chan events.DataEvent, 1)
	f.rowerDataEvents = &dataEvents

	for {
		select {
		case <-n.Context().Done():
			log.Println("Subscription stopped")
			f.rowerDataEvents = nil
			return

		case dataEvent := <-dataEvents:

			featureData := make([]byte, 0)
			// ?   0 .. Stroke rate and Stroke count (1 if NOT present)
			// 0   1 .. Average Stroke rate (1 if present)
			// ?   2 .. Total Distance present (1 if present)
			// ?   3 .. Instantaneous Pace (1 if present)
			// 0   4 .. Average Pace (1 if present)
			// ?   5 .. Instantaneous Power (1 if present)
			// 0   6 .. Average Power (1 if present)
			// 0   7 .. Resistance Level (1 if present)
			featureData = append(featureData, byte(0x01))
			// ?   8 .. Expended Energy (1 if present)
			// 0   9 .. Heart Rate (1 if present)
			// 0  10 .. Metabolic Equivalent (1 if present)
			// ?  11 .. Elapsed Time in seconds (1 if present)
			// 0  12 .. Remaining Time (1 if present)
			// 0  13 .. Reserved for future use
			// 0  14 .. Reserved for future use
			// 0  15 .. Reserved for future use
			featureData = append(featureData, byte(0x00))

			// Bit 0 - Stroke rate and Stroke count (1 if NOT present)
			if dataEvent.StrokesPerMinute != nil && dataEvent.Strokes != nil {
				strokeRate := uint8(*(dataEvent.StrokesPerMinute) * 2)
				featureData = append(featureData, strokeRate)
				strokeCount := make([]byte, 2)
				binary.LittleEndian.PutUint16(strokeCount, *dataEvent.Strokes)
				featureData = append(featureData, strokeCount...)
				featureData[0] ^= 1 << 0
			}

			// Bit 2 - Total Distance
			if dataEvent.Distance != nil {
				distance := *dataEvent.Distance
				totalDistance := make([]byte, 3)
				totalDistance[0] = byte(distance & 255)
				totalDistance[1] = byte((distance >> 8) & 255)
				totalDistance[2] = 0
				featureData = append(featureData, totalDistance...)
				featureData[0] |= 4
			}

			// Bit 3 - Instantaneous Pace
			if dataEvent.Time500mSplit != nil {
				instantaneousPace := make([]byte, 2)
				binary.LittleEndian.PutUint16(instantaneousPace, *dataEvent.Time500mSplit)
				featureData = append(featureData, instantaneousPace...)
				featureData[0] |= 8
			}

			// Bit 5 - Instantaneous Power
			if dataEvent.WattsPreviousStroke != nil {
				instantaneousPower := make([]byte, 2)
				binary.LittleEndian.PutUint16(instantaneousPower, *dataEvent.WattsPreviousStroke)
				featureData = append(featureData, instantaneousPower...)
				featureData[0] |= 32
			}

			// Bit 11 - Elapsed Time in seconds
			elapsedTime := make([]byte, 2)
			binary.LittleEndian.PutUint16(elapsedTime, dataEvent.ElapsedTime)
			featureData = append(featureData, elapsedTime...)
			featureData[1] |= 8

			_, err := n.Write(featureData)
			if err != nil {
				log.Printf("Error writing feature data: %+v\n", err)
			}

		case <-time.After(time.Minute * 5):
			log.Println("Timeout")
		}
	}
}

func (f *FitnessMachine) rowerFeatureReadHandler(_ ble.Request, rsp ble.ResponseWriter) {
	log.Println("Rower feature read request.")

	// ?   0 .. Stroke rate and Stroke count (1 if NOT present)
	// 0   1 .. Average Stroke rate (1 if present)
	// 1   2 .. Total Distance present
	// ?   3 .. Instantaneous Pace (1 if present)
	// 0   4 .. Average Pace (1 if present)
	// ?   5 .. Instantaneous Power (1 if present)
	// 0   6 .. Average Power (1 if present)
	// 0   7 .. Resistance Level (1 if present)

	// ?   8 .. Expended Energy (1 if present)
	// 0   9 .. Heart Rate (1 if present)
	// 0  10 .. Metabolic Equivalent (1 if present)
	// 1  11 .. Elapsed Time in seconds (1 if present)
	// 0  12 .. Remaining Time (1 if present)
	// 0  13 .. Reserved for future use
	// 0  14 .. Reserved for future use
	// 0  15 .. Reserved for future use

	_, err := rsp.Write([]byte{0x2C, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		log.Println("Rower feature read request failed")
	}
}

func (f *FitnessMachine) OnData(event events.DataEvent) {
	if f.rowerDataEvents != nil {
		*f.rowerDataEvents <- event
	}
}

func (f *FitnessMachine) OnStatusChange(event events.StatusChangeEvent) {
	log.Printf("OnStatusChangeEvent: %+v\n", event)
}
