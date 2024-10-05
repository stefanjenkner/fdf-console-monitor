package mocks

import (
	"bytes"
	"fmt"
	"time"

	"go.bug.st/serial"
)

type MockSerialPort struct {
	readBuffer  bytes.Buffer
	writeBuffer bytes.Buffer
	Closed      int
}

func NewMockSerialPort(bufferString *bytes.Buffer) (*MockSerialPort, serial.Port) {
	mockSerialPort := &MockSerialPort{
		readBuffer:  *bufferString,
		writeBuffer: bytes.Buffer{},
	}
	var port serial.Port = mockSerialPort
	return mockSerialPort, port
}

func (m *MockSerialPort) Read(p []byte) (n int, err error) {
	read, err := m.readBuffer.Read(p)
	if err != nil {
		return read, fmt.Errorf("read: unexpected error: %w", err)
	}
	return read, nil
}

func (m *MockSerialPort) Write(p []byte) (n int, err error) {
	write, err := m.writeBuffer.Write(p)
	if err != nil {
		return 0, fmt.Errorf("write: unexpected error: %w", err)
	}
	return write, nil
}

func (m *MockSerialPort) Close() error {
	m.Closed++
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
