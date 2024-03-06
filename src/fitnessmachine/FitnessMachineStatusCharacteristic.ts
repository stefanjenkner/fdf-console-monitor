import { Characteristic } from '@abandonware/bleno'
import { StatusChange } from '../monitor/StatusChange';
import { UpdateValueCallback } from './UpdateValueCallback';
import log from 'loglevel'

export class FitnessMachineStatusCharacteristic extends Characteristic {

    updateValueCallback?: UpdateValueCallback | null;
    maxValueSize?: number | null;

    constructor() {
        super({
            uuid: '2ADA',
            value: null,
            properties: ['notify']
        });
    }

    onSubscribe(maxValueSize: number, updateValueCallback: UpdateValueCallback): void {
        log.debug(`FitnessMachineStatusCharacteristic onSubscribe maxValueSize=${maxValueSize}`)
        this.updateValueCallback = updateValueCallback;
        this.maxValueSize = maxValueSize;
    }

    onUnsubscribe(): void {
        log.debug('FitnessMachineStatusCharacteristic onUnsubscribe')
        this.updateValueCallback = null;
        this.maxValueSize = null;
    }

    onStatusChange(statusChange: StatusChange): void {
        const result = Buffer.alloc(2);
        switch (statusChange) {
            case StatusChange.Started:
            case StatusChange.Resumed:
                // 0x04 - Fitness Machine Started or Resumed by the User
                result.writeUInt8(0x04, 0);
                break;
            case StatusChange.PausedOrStopped:
                // 0x02 - Fitness Machine Stopped or Paused by the User
                result.writeUInt8(0x02, 0);
                break;
            case StatusChange.Reset:
                // 0x01 - Reset
                result.writeUInt8(0x01, 0);
                break;
            default:
                log.warn(`Unsupported status change: ${statusChange}`)
                return;
        }
        this.updateValueCallback && this.updateValueCallback(result);
    }
}