package serialmonitor

import (
	"bufio"
	"fmt"
	"log"
	"strings"

	"github.com/stefanjenkner/fdf-console-monitor/internal/events"
	"github.com/stefanjenkner/fdf-console-monitor/internal/observer"
	"go.bug.st/serial"
)

type SerialMonitor struct {
	portName    string
	port        *serial.Port
	observers   map[observer.Observer]struct{}
	stopChannel *chan struct{}
}

func NewSerialMonitor(portName string) *SerialMonitor {
	return &SerialMonitor{
		portName:  portName,
		observers: make(map[observer.Observer]struct{}),
	}
}

func (m *SerialMonitor) Run() {
	log.Printf("Running SerialMonitor on port: %s\n", m.portName)

	stopChan := make(chan struct{})
	m.stopChannel = &stopChan
	defer close(stopChan)

	strokes := uint16(0)
	isPausedOrStopped := false
	lineChannel := m.createLineChannel(&stopChan)

	for line := range lineChannel {
		log.Printf("Received: %s\n", line)
		switch {
		case strings.HasPrefix(line, "A"):
			capture := parse(line)
			if capture.strokesPerMinute == 0 {
				m.emitStatusChangeEvent(*events.NewStatusChangeEvent(events.PausedOrStopped))
				isPausedOrStopped = true
			} else if isPausedOrStopped {
				m.emitStatusChangeEvent(*events.NewStatusChangeEvent(events.Resumed))
				isPausedOrStopped = false
			} else if strokes == 0 {
				m.emitStatusChangeEvent(*events.NewStatusChangeEvent(events.Started))
			}
			if isPausedOrStopped {
				m.emitDataEvent(*events.NewDataEvent(capture.elapsedTime, capture.level,
					events.WithRemainingDistance(capture.distance),
					events.WithTime500mAverage(capture.time500m),
					events.WithWattsAverage(capture.watts),
					events.WithCaloriesTotal(capture.cals),
				))
			} else {
				strokes++
				m.emitDataEvent(*events.NewDataEvent(capture.elapsedTime, capture.level,
					events.WithDistance(capture.distance),
					events.WithStrokes(strokes),
					events.WithStrokesPerMinute(capture.strokesPerMinute),
					events.WithTime500mSplit(capture.time500m),
					events.WithWattsPreviousStroke(capture.watts),
					events.WithCaloriesPerHour(capture.cals),
				))
			}

		case strings.HasPrefix(line, "W"):
			if err := m.writeLine("K"); err != nil {
				log.Println(err)
			}

		case strings.HasPrefix(line, "R"):
			m.emitStatusChangeEvent(*events.NewStatusChangeEvent(events.Reset))
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

func (m *SerialMonitor) AddObserver(o observer.Observer) {
	m.observers[o] = struct{}{}
}

func (m *SerialMonitor) emitDataEvent(event events.DataEvent) {
	for o := range m.observers {
		o.OnData(event)
	}
}

func (m *SerialMonitor) emitStatusChangeEvent(event events.StatusChangeEvent) {
	for o := range m.observers {
		o.OnStatusChange(event)
	}
}

func (m *SerialMonitor) writeLine(line string) error {
	_, err := (*m.port).Write([]byte(line + "\n"))
	if err != nil {
		return fmt.Errorf("writeLine: unexpected error: %w", err)
	}
	return nil
}

func (m *SerialMonitor) closePort() {
	if m.port == nil {
		return
	}
	if err := (*m.port).Close(); err != nil {
		log.Printf("Error closing port: %+v\n", err)
	}
}
