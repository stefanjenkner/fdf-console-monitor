#!/usr/bin/env node

import { FitnessMachine } from './fitnessmachine/FitnessMachine'
import { Monitor } from './monitor/Monitor';
import log from 'loglevel'
import { parseArgs } from 'node:util';

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
        shutdown(1);
        return;
    }
    fitnessMachine.start();
});
monitor.on('disconnect', (error?) => {
    if (error) {
        shutdown(1);
        return;
    }
    fitnessMachine.stop();
});
monitor.on('data', (data) => fitnessMachine.onData(data));

process.on('SIGINT', function () {
    shutdown(0);
});

monitor.connect();

function shutdown(exitCode : number) {
    monitor.disconnect();
    fitnessMachine.stop();
    setTimeout(() => {
        log.info('Bye Bye');
        process.exit(exitCode);
    }, 3000);
}
