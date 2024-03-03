import log from 'loglevel'

import { SerialPort } from 'serialport'
import { ReadlineParser } from '@serialport/parser-readline'
import { Capture } from './Capture';
import { Parser } from './Parser';

interface MonitorOptions {
    port: string
}

export type ConnectCallback = (err: Error | null) => void;
export type DisconnectCallback = (err: Error | null) => void;
export type OnDataCallback = (capture: Capture) => void;

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
            const captureParser = new Parser();
            parser.on('data', (data : string) => {

                log.debug(`Received: ${data}`);

                if (data.startsWith('A')) {
                    this.onDataCallback && this.onDataCallback(captureParser.parse(data));
                } else if (data.startsWith('W')) {
                    port.write('K\n')
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