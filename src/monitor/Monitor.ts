import { Data } from './Data';
import EventEmitter from 'events';
import { Parser } from './Parser';
import { ReadlineParser } from '@serialport/parser-readline'
import { SerialPort } from 'serialport'
import TypedEmitter from 'typed-emitter'
import log from 'loglevel'

interface MonitorOptions {
    port: string
}

type MonitorEvents = {
    'connect': (err: Error | null) => void,
    'disconnect': (err: Error | null) => void,
    'data': (data: Data) => void
}

export class Monitor extends(EventEmitter as new () => TypedEmitter<MonitorEvents>) {

    serialPort?: SerialPort;
    options: MonitorOptions;

    constructor(options: MonitorOptions) {
        super();
        this.options = options;
    }

    connect(): void {
        const port = new SerialPort({
            path: this.options.port,
            baudRate: 9600,
        }, (error?) => {
            if (error) {
                log.error(`Error opening serial port: ${error.message}`);
                this.emit('connect', error);
                return;
            }

            port.write('C\n');
            const parser = port.pipe(new ReadlineParser());
            let strokes = 0;
            const captureParser = new Parser();
            parser.on('data', (rawData : string) => {

                log.debug(`Received: ${rawData}`);

                if (rawData.startsWith('A')) {
                    const capture = captureParser.parse(rawData);
                    const isPausedOrStopped = capture.strokesPerMinute === 0;
                    if (!isPausedOrStopped) strokes++;
                    const data : Data = {
                        elapsedTime: capture.elapsedTime,
                        distance: capture.distance,
                        strokes: isPausedOrStopped ? null : strokes,
                        strokesPerMinute: isPausedOrStopped ? null : capture.strokesPerMinute,
                        level: capture.level,
                        time500mSplit: isPausedOrStopped ? null : capture.time500m,
                        time500mAverage: isPausedOrStopped ? capture.time500m : null,
                        wattsPreviousStroke: isPausedOrStopped ? null : capture.watts,
                        wattsAverage: isPausedOrStopped ? capture.watts : null,
                        caloriesPerHour: isPausedOrStopped ? null : capture.cals,
                        caloriesTotal: isPausedOrStopped ? capture.cals : null,
                    }
                    this.emit('data', data);
                } else if (rawData.startsWith('W')) {
                    port.write('K\n')
                } else if (rawData.startsWith('R')) {
                    strokes = 0;
                }
            });

            log.info('Connection established.');
            this.serialPort = port;
            this.emit('connect', null);
        })
    }

    disconnect(): void {

        if (this.serialPort) {
            this.serialPort.removeAllListeners();
            this.serialPort.writable && this.serialPort.write('D\n');
            if (!(this.serialPort.closed || this.serialPort.closing)) {
                log.info(`Closing serial port: ${this.options.port}`)
                this.serialPort.close((error) => {
                    if (error) {
                        log.error(`Disconnect error: ${error.message}`);
                        this.emit('disconnect', error);
                        return;
                    }
                });
            }
        }

        log.info('Connection closed.');
        this.emit('disconnect', null);
    }
}