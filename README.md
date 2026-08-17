# FDF Console monitor

Expose a First Degree Fitness water rower as BLE (Bluetooth Low Energy)
peripheral providing Fitness Machine Service (FTMS) rower data.

Tested with:

- First Degree Fitness NEON plus water rower which comes with the (basic) FDF
  Console and a serial interface.
- First Degree Fitness Viking 2 AR — community-reported, works unmodified
  ([#168](https://github.com/stefanjenkner/fdf-console-monitor/issues/168)).

## Acknowledgements

Thanks to [@evelynDickinson](https://github.com/evelynDickinson) for the
Viking 2 AR hardware report and protocol observations (`C<nnnn>` handshake ACK
payload, `H<nnn>` heart-rate lines).

## Usage

Run:

```bash
go run ./cmd/fdf-console-monitor --name "FDF Rower" --port /dev/ttyUSB0
```

## Advanced usage: building on Linux for running in non-root context

Build binary on Linux - e.g. for running on Raspberry Pi:

```bash
CGO_ENABLED=0 go build -a -o fdf-console-monitor ./cmd/fdf-console-monitor
```

Grant only specific capabilities instead of full root access:

```bash
sudo setcap 'cap_net_raw,cap_net_admin+eip' ./fdf-console-monitor
```

Run as non-root:

```bash
./fdf-console-monitor --name "FDF Rower" --port /dev/ttyUSB0
```

## Contribution

Prerequisites for development:

<details>
<summary>macOS</summary>

```bash
brew install pre-commit commitizen golangci-lint
```

</details>

Set up pre-commit hooks:

```bash
pre-commit install && pre-commit install --hook-type commit-msg && pre-commit run
```

## Notes

Uses [go-ble](https://github.com/go-ble/ble) for BLE communication
