# FDF Console monitor

Expose a First Degree Fitness water rower as BLE (Bluetooth Low Energy) peripheral providing Fitness Machine Service (FTMS) rower data.

Tested with First Degree Fitness NEON plus water rower which comes with the (basic) FDF Console and a serial interface.

## Usage

Install dependencies and build:

    npm install && npm run build

Run:

    npx fdf-console-monitor --port /dev/ttyUSB0

## Notes

Uses [bleno](bleno) for BLE communication - please check [prerequisites] and hints for [running on Linux]

[bleno]: https://github.com/abandonware/bleno
[prerequisites]: https://github.com/abandonware/bleno?tab=readme-ov-file#prerequisites
[running on Linux]: https://github.com/abandonware/bleno?tab=readme-ov-file#running-on-linux