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

        log.debug(`RowerDataCharacteristic onUnsubscribe`)
        this.updateValueCallback = null;
        this.maxValueSize = null;  
    }

    onCapture(capture : Capture): void {

        const flags1 = Buffer.alloc(1);
        // *  0 .. More Data (inverted)
        // *  1 .. Average Stroke rate present (inverted)
        // *  2 .. Total Distance present
        //    3 .. Instantaneous Pace present
        //    4 .. Average Pace present
        //    5 .. Instantaneous Power present
        //    6 .. Average Power present
        //    7 .. Resistance Level present
        flags1.writeUInt8(0x07);

        const flags2 = Buffer.alloc(1);
        //    8 .. Expended Energy present
        //    9 .. Heart Rate present
        //   10 .. Metabolic Equivalent present
        // * 11 .. Elapsed Time present
        //   12 .. Remaining Time present
        //   13 .. Reserved for future use
        //   14 .. Reserved for future use
        //   15 .. Reserved for future use
        flags2.writeUInt8(0x08);

        const totalDistance = Buffer.alloc(3);
        totalDistance.writeUInt8((capture.distance || 0) & 255)
        totalDistance.writeUInt16LE((capture.distance || 0) >> 8, 1)
        
        const elapsedTime = Buffer.alloc(2);
        elapsedTime.writeUInt16LE(capture.elapsedTime || 0)
        
        this.updateValueCallback && this.updateValueCallback(Buffer.concat([flags1, flags2, totalDistance, elapsedTime]));
    }
}