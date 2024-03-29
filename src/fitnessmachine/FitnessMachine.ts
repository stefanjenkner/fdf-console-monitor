import { Data } from '../Data';
import { FitnessMachineOptions } from './FitnessMachineOptions';
import { FitnessMachineService } from './FitnessMachineService';
import { StatusChange } from '../StatusChange';
import bleno from '@abandonware/bleno';
import log from 'loglevel';

export class FitnessMachine {

    options: FitnessMachineOptions;
    fitnessMachineService: FitnessMachineService;

    constructor(options: FitnessMachineOptions) {
        this.options = options;
        this.fitnessMachineService = new FitnessMachineService();
    }

    start(): void {

        bleno.on('stateChange', (state) => {
            log.info(`State changed to: ${state}`);
            if (state === 'poweredOn') {
                bleno.startAdvertising(this.options.name, [this.fitnessMachineService.uuid], (error) => {
                    if (error) {
                        log.error(error);
                    }
                });
            } else {
                bleno.stopAdvertising();
            }
        });

        bleno.on('accept', (clientAddress) => {
            log.debug(`connected: ${clientAddress}`);
            bleno.updateRssi();
        });

        bleno.on('disconnect', (clientAddress) => {
            log.debug(`disconnected: ${clientAddress}`);
        });

        bleno.on('advertisingStart', (error) => {
            if (!error) {
                bleno.setServices(
                    [this.fitnessMachineService],
                    (error) => {
                        if (error) {
                            log.error(`Set service error: ${error}`);
                        }
                    });
            }
        });

        bleno.on('advertisingStop', () => {
            log.info('Advertising stopped');
        });
    }

    stop(): void {
        bleno.removeAllListeners();
        bleno.disconnect();
        bleno.stopAdvertising(() => {
            log.info('Advertising stopped successfully');
        });
    }

    onData(data : Data) {
        this.fitnessMachineService.onData(data);
    }

    onStatusChange(statusChange : StatusChange) {
        this.fitnessMachineService.onStatusChange(statusChange);
    }
}
