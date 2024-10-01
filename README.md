# FDF Console monitor

Expose a First Degree Fitness water rower as BLE (Bluetooth Low Energy) peripheral providing Fitness Machine Service
(FTMS) rower data.

Tested with First Degree Fitness NEON plus water rower which comes with the (basic) FDF Console and a serial interface.

## Usage

Install dependencies and build:

    go mod download

Run:

    go run ./cmd/app --name "FDF Rower" --port /dev/ttyUSB0

Optional: Build binary for Linux and set capability flags:

    CGO_ENABLED=0 go build -a -o fdf-console-monitor ./cmd/app
    sudo setcap 'cap_net_raw,cap_net_admin+eip' fdf-console-monitor
    ./fdf-console-monitor --name "FDF Rower" --port /dev/ttyUSB0

## Notes

Uses [go-ble](https://github.com/go-ble/ble) for BLE communication
