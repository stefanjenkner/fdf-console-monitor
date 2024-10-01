package serialmonitor

import (
	"bufio"
	"fdf-console-monitor/internal/events"
	"fmt"
	"go.bug.st/serial"
	"log"
	"strings"
)

type SerialMonitor struct {
	portName    string
	port        *serial.Port
	observers   map[events.Observer]struct{}
	stopChannel *chan struct{}
}

func NewSerialMonitor(portName string) *SerialMonitor {
	return &SerialMonitor{
		portName:  portName,
		observers: make(map[events.Observer]struct{}),
	}
}

func (m *SerialMonitor) Run() {
	log.Printf("Running SerialMonitor on port: %s\n", m.portName)

	stopChan := make(chan struct{})
	m.stopChannel = &stopChan
	defer close(stopChan)

	strokes := uint64(0)
	isPausedOrStopped := false
	lineChannel := m.createLineChannel(&stopChan)

	for line := range lineChannel {

		log.Printf("Received: %s\n", line)
		switch {
		case strings.HasPrefix(line, "A"):
			capture := parse(line)
			if capture.strokesPerMinute == 0 {
				m.emitStatusChangeEvent(events.StatusChangeEvent{StatusChange: events.PausedOrStopped})
				isPausedOrStopped = true
			} else if isPausedOrStopped {
				m.emitStatusChangeEvent(events.StatusChangeEvent{StatusChange: events.Resumed})
				isPausedOrStopped = false
			} else if strokes == 0 {
				m.emitStatusChangeEvent(events.StatusChangeEvent{StatusChange: events.Started})
			}
			builder := events.NewDataEventBuilder(capture.elapsedTime, capture.level)
			if isPausedOrStopped {
				builder.SetRemainingDistance(capture.distance)
				builder.SetTime500mAverage(capture.time500m)
				builder.SetWattsAverage(capture.watts)
				builder.SetCaloriesTotal(capture.cals)
			} else {
				strokes++
				builder.SetDistance(capture.distance)
				builder.SetStrokes(strokes)
				builder.SetStrokesPerMinute(capture.strokesPerMinute)
				builder.SetTime500mSplit(capture.time500m)
				builder.SetWattsPreviousStroke(capture.watts)
				builder.SetCaloriesPerHour(capture.cals)
			}
			m.emitDataEvent(*builder.Build())

		case strings.HasPrefix(line, "W"):
			if err := m.writeLine("K"); err != nil {
				log.Println(err)
			}

		case strings.HasPrefix(line, "R"):
			m.emitStatusChangeEvent(events.StatusChangeEvent{StatusChange: events.Reset})
			isPausedOrStopped = false
			strokes = 0
		}
	}

	log.Println("Stopped SerialMonitor")
}

func (m *SerialMonitor) createLineChannel(stop *chan struct{}) <-chan string {
	channel := make(chan string)

	go func() {
		defer close(channel)

		// connect
		if err := m.openPort(); err != nil {
			log.Printf("Error opening port: %s\n", err)
			return
		}
		defer m.closePort()

		// send C for connect
		if err := m.writeLine("C"); err != nil {
			log.Printf("Error connecting: %s\n", err)
			return
		}

		// read line by line until EOT or receiving stop
		scanner := bufio.NewScanner(*m.port)
		for scanner.Scan() {
			select {
			case channel <- scanner.Text():
			case <-*stop:
				log.Println("SerialMonitor received stop signal")
				return
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("SerialMonitor received non-EOF error: %s\n", err)
			return
		}
	}()

	return channel
}

func (m *SerialMonitor) openPort() error {
	if m.port == nil {
		mode := &serial.Mode{BaudRate: 9600}
		port, err := serial.Open(m.portName, mode)
		if err != nil {
			return fmt.Errorf("failed to open port %s: %w", m.portName, err)
		}
		m.port = &port
	}
	return nil
}

func (m *SerialMonitor) Stop() {
	if (*m.stopChannel) != nil {
		log.Println("Stopping SerialMonitor...")
		*m.stopChannel <- struct{}{}
	}
}

func (m *SerialMonitor) AddObserver(observer events.Observer) {
	m.observers[observer] = struct{}{}
}

func (m *SerialMonitor) emitDataEvent(event events.DataEvent) {
	for observer := range m.observers {
		observer.OnData(event)
	}
}

func (m *SerialMonitor) emitStatusChangeEvent(event events.StatusChangeEvent) {
	for observer := range m.observers {
		observer.OnStatusChange(event)
	}
}

func (m *SerialMonitor) writeLine(line string) error {
	_, err := (*m.port).Write([]byte(line + "\n"))
	return err
}

func (m *SerialMonitor) closePort() {
	if m.port == nil {
		return
	}
	if err := (*m.port).Close(); err != nil {
		log.Printf("Error closing port: %+v\n", err)
	}
}
