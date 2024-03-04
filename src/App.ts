#!/usr/bin/env node

import { parseArgs } from 'node:util';
import log from 'loglevel'
import { FitnessMachine } from './fitnessmachine/FitnessMachine'
import { Monitor } from './monitor/Monitor';
import { Data } from './monitor/Data';

log.setLevel('DEBUG')

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

const { values: { name, port } } = parseArgs({ options });

const fitnessMachine = new FitnessMachine({ name: name ? name : 'FDF Rower' })
const monitor = new Monitor({ port: port ? port : '/dev/ttyUSB0' }, (data: Data) => {
    fitnessMachine.onData(data);
});
monitor.connect((error?) => {
    if (error) {
        process.exit(1);
    }
    fitnessMachine.start();
});

process.on('SIGINT', function () {
    let exitCode = 1;
    monitor.disconnect((error?) => {
        if (!error) {
            exitCode = 0;
        }
    })
    fitnessMachine.stop();

    setTimeout(() => {
        log.info('Bye Bye');
        process.exit(exitCode);
    }, 3000);
});
