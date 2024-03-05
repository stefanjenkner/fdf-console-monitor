#!/usr/bin/env node

import { FitnessMachine } from './fitnessmachine/FitnessMachine'
import { Monitor } from './monitor/Monitor';
import log from 'loglevel'
import { parseArgs } from 'node:util';
import process from 'node:process';

const defaultPort = '/dev/ttyUSB0';
const defaultName = 'FDF Rower';

const options = {
    name: {
        type: 'string',
        short: 'n'
    },
    port: {
        type: 'string',
        short: 'p'
    }
} as const

log.setLevel('DEBUG')
const { values: { name, port } } = parseArgs({ options });
const fitnessMachine = new FitnessMachine({ name: name ? name : defaultName })
const monitor = new Monitor({ port: port ? port : defaultPort });
monitor.on('connect', (error?) => {
    if (error) {
        process.exitCode = 1;
    } else {
        fitnessMachine.start();
    }
});
monitor.on('disconnect', (error?) => {
    fitnessMachine.stop();
    process.exitCode = error ? 1 : 0
});
monitor.on('data', (data) => fitnessMachine.onData(data));

process.on('SIGINT', function () {
    monitor.disconnect();
    fitnessMachine.stop();
    process.exitCode = 0
});

monitor.connect();
