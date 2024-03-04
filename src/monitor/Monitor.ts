import log from 'loglevel'

import { SerialPort } from 'serialport'
import { ReadlineParser } from '@serialport/parser-readline'
import { Parser } from './Parser';
import { Data } from './Data';

interface MonitorOptions {
    port: string
}

export type ConnectCallback = (err: Error | null) => void;
export type DisconnectCallback = (err: Error | null) => void;
export type OnDataCallback = (data: Data) => void;

export class Monitor {

    serialPort?: SerialPort;
    options: MonitorOptions;
    onDataCallback: OnDataCallback;

    constructor(options: MonitorOptions, onDataCallback: OnDataCallback) {
        this.options = options;
        this.onDataCallback = onDataCallback;
    }

    connect(callback?: ConnectCallback): void {
        const port = new SerialPort({
            path: this.options.port,
            baudRate: 9600,
        }, (error?) => {
            if (error) {
                log.error(`Error opening serial port: ${error.message}`);
                callback && callback(error || null);
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
                    this.onDataCallback && this.onDataCallback(data);
                } else if (rawData.startsWith('W')) {
                    port.write('K\n')
                } else if (rawData.startsWith('R')) {
                    strokes = 0;
                }
            });

            log.info('Connection established.');
            this.serialPort = port;
            callback && callback(null)
        })
    }

    disconnect(callback?: DisconnectCallback): void {

        if (this.serialPort) {
            this.serialPort.removeAllListeners();
            this.serialPort.writable && this.serialPort.write('D\n');
            if (!(this.serialPort.closed || this.serialPort.closing)) {
                log.info(`Closing serial port: ${this.options.port}`)
                this.serialPort.close((error) => {
                    if (error) {
                        log.error(`Disconnect error: ${error.message}`);
                        callback && callback(error || null);
                        return;
                    }
                });
            }
        }

        log.info('Connection closed.');
        callback && callback(null);
    }
}