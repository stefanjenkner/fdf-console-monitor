import log from 'loglevel'
import { Characteristic } from '@abandonware/bleno'

type ReadRequestCallback = (result: number, data?: Buffer) => void;

export class RowerFeatureCharacteristic extends Characteristic {

    constructor() {
        super({
            uuid: '2ACC',
            properties: ['read'],
            value: null
        })
    }

    onReadRequest(offset: number, callback: ReadRequestCallback): void {

        log.debug(`RowerFeatureCharacteristic onReadRequest offset=${offset}`);

        // 1   0 .. Stroke rate and Stroke count (1 if NOT present)
        // 0   1 .. Average Stroke rate (1 if present)
        // 1   2 .. Total Distance present
        // 1   3 .. Instantaneous Pace (1 if present)
        // 0   4 .. Average Pace (1 if present)
        // 1   5 .. Instantaneous Power (1 if present)
        // 0   6 .. Average Power (1 if present)
        // 0   7 .. Resistance Level (1 if present)

        // 1   8 .. Expended Energy (1 if present)
        // 0   9 .. Heart Rate (1 if present)
        // 0  10 .. Metabolic Equivalent (1 if present)
        // 1  11 .. Elapsed Time in seconds (1 if present)
        // 0  12 .. Remaining Time (1 if present)
        // 0  13 .. Reserved for future use
        // 0  14 .. Reserved for future use
        // 0  15 .. Reserved for future use

        const features = [0x2D, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        callback(this.RESULT_SUCCESS, Buffer.from(features.slice(offset, features.length)));
    }
}