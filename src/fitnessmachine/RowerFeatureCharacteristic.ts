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

        // *  0 .. More Data (inverted)
        // *  1 .. Average Stroke rate present (inverted)
        // *  2 .. Total Distance present
        //    3 .. Instantaneous Pace present
        //    4 .. Average Pace present
        //    5 .. Instantaneous Power present
        //    6 .. Average Power present
        //    7 .. Resistance Level present

        //    8 .. Expended Energy present
        //    9 .. Heart Rate present
        //   10 .. Metabolic Equivalent present
        // * 11 .. Elapsed Time present
        //   12 .. Remaining Time present
        //   13 .. Reserved for future use
        //   14 .. Reserved for future use
        //   15 .. Reserved for future use

        const features = [0x07, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00];
        callback(this.RESULT_SUCCESS, Buffer.from(features.slice(offset, features.length)));
    }
}