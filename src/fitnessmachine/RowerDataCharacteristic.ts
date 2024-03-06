import { Characteristic } from '@abandonware/bleno'
import { Data } from '../monitor/Data';
import { UpdateValueCallback } from './UpdateValueCallback';
import log from 'loglevel'

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

    onSubscribe(maxValueSize: number, updateValueCallback: UpdateValueCallback): void {

        log.debug(`RowerDataCharacteristic onSubscribe maxValueSize=${maxValueSize}`)
        this.updateValueCallback = updateValueCallback;
        this.maxValueSize = maxValueSize;
    }

    onUnsubscribe(): void {

        log.debug('RowerDataCharacteristic onUnsubscribe')
        this.updateValueCallback = null;
        this.maxValueSize = null;
    }

    onData(data: Data): void {

        const featureData: Array<Buffer> = []
        // ?   0 .. Stroke rate and Stroke count (1 if NOT present)
        // 0   1 .. Average Stroke rate (1 if present)
        // 1   2 .. Total Distance present
        // ?   3 .. Instantaneous Pace (1 if present)
        // 0   4 .. Average Pace (1 if present)
        // ?   5 .. Instantaneous Power (1 if present)
        // 0   6 .. Average Power (1 if present)
        // 0   7 .. Resistance Level (1 if present)
        let featuresOctet1 = 0x05;
        // ?   8 .. Expended Energy (1 if present)
        // 0   9 .. Heart Rate (1 if present)
        // 0  10 .. Metabolic Equivalent (1 if present)
        // 1  11 .. Elapsed Time in seconds (1 if present)
        // 0  12 .. Remaining Time (1 if present)
        // 0  13 .. Reserved for future use
        // 0  14 .. Reserved for future use
        // 0  15 .. Reserved for future use
        const featuresOctet2 = 0x08;

        // Bit 0 - Stroke rate and Stroke count (1 if NOT present)
        if (data.strokesPerMinute && data.strokes) {
            const stokeRate = Buffer.alloc(1);
            stokeRate.writeUInt8(Math.round(data.strokesPerMinute * 2) || 0);
            const strokeCount = Buffer.alloc(2);
            strokeCount.writeUInt16LE(data.strokes || 0)
            featuresOctet1 &= ~1;
            featureData.push(stokeRate);
            featureData.push(strokeCount);
        }

        // Bit 2 - Total Distance
        const totalDistance = Buffer.alloc(3);
        totalDistance.writeUInt8((data.distance || 0) & 255)
        totalDistance.writeUInt16LE((data.distance || 0) >> 8, 1)
        featureData.push(totalDistance);

        // Bit 3 - Instantaneous Pace
        if (data.time500mSplit) {
            const instantaneousPace = Buffer.alloc(2);
            instantaneousPace.writeUInt16LE(data.time500mSplit || 0)
            featuresOctet1 |= 8;
            featureData.push(instantaneousPace);
        }

        // Bit 5 - Instantaneous Power
        if (data.wattsPreviousStroke) {
            const instantaneousPower = Buffer.alloc(2);
            instantaneousPower.writeUInt16LE(data.wattsPreviousStroke || 0)
            featuresOctet1 |= 32;
            featureData.push(instantaneousPower);
        }

        // Bit 8 - Expended Energy
        // const totalEnergy = Buffer.alloc(2);
        // totalEnergy.writeUInt16LE(0)
        // const energyPerHour = Buffer.alloc(2);
        // energyPerHour.writeUInt16LE(data.caloriesPerHour || 0)
        // const energyPerMinute = Buffer.alloc(1);
        // energyPerMinute.writeUInt8(0)

        // Bit 11 - Elapsed Time in seconds
        const elapsedTime = Buffer.alloc(2);
        elapsedTime.writeUInt16LE(data.elapsedTime || 0)
        featureData.push(elapsedTime);

        // Feature flags
        const featureFlags = Buffer.alloc(2);
        featureFlags.writeUInt8(featuresOctet1 || 0);
        featureFlags.writeUInt8(featuresOctet2 || 0, 1);
        featureData.unshift(featureFlags);

        this.updateValueCallback && this.updateValueCallback(Buffer.concat(featureData));
    }
}