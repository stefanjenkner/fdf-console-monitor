import log from 'loglevel'
import { Characteristic } from '@abandonware/bleno'
import { Capture } from '../monitor/Capture';

type UpdateValueCallback = (data: Buffer) => void;

export class RowerDataCharacteristic extends Characteristic {

    updateValueCallback?: UpdateValueCallback | null;
    maxValueSize?: number | null;

    constructor() {
        super({
            uuid: '2AD1',
            value: null,
            properties: ['notify']
        });
        this.updateValueCallback = null;
        this.maxValueSize = null;
    }

    onSubscribe(maxValueSize: number, updateValueCallback: (data: Buffer) => void): void {

        log.debug(`RowerDataCharacteristic onSubscribe maxValueSize=${maxValueSize}`)
        this.updateValueCallback = updateValueCallback;
        this.maxValueSize = maxValueSize;
    }

    onUnsubscribe(): void {

        log.debug('RowerDataCharacteristic onUnsubscribe')
        this.updateValueCallback = null;
        this.maxValueSize = null;
    }

    onCapture(capture: Capture): void {

        const flags = Buffer.alloc(2);
        // 1   0 .. Stroke rate and Stroke count (1 if NOT present)
        // 0   1 .. Average Stroke rate (1 if present)
        // 1   2 .. Total Distance present
        // 1   3 .. Instantaneous Pace (1 if present)
        // 0   4 .. Average Pace (1 if present)
        // 0   5 .. Instantaneous Power (1 if present)
        // 0   6 .. Average Power (1 if present)
        // 0   7 .. Resistance Level (1 if present)
        flags.writeUInt8(0x0D || 0);
        // 0   8 .. Expended Energy (1 if present)
        // 0   9 .. Heart Rate (1 if present)
        // 0  10 .. Metabolic Equivalent (1 if present)
        // 1  11 .. Elapsed Time in seconds (1 if present)
        // 0  12 .. Remaining Time (1 if present)
        // 0  13 .. Reserved for future use
        // 0  14 .. Reserved for future use
        // 0  15 .. Reserved for future use
        flags.writeUInt8(0x08 || 0, 1);

        const totalDistance = Buffer.alloc(3);
        totalDistance.writeUInt8((capture.distance || 0) & 255)
        totalDistance.writeUInt16LE((capture.distance || 0) >> 8, 1)

        const instantaneousPace = Buffer.alloc(2);
        instantaneousPace.writeUInt16LE(capture.time500mSplit || 0)

        const elapsedTime = Buffer.alloc(2);
        elapsedTime.writeUInt16LE(capture.elapsedTime || 0)

        this.updateValueCallback && this.updateValueCallback(Buffer.concat([flags, totalDistance, instantaneousPace, elapsedTime]));
    }
}