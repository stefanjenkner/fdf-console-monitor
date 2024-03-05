
import log from 'loglevel'

export class RowerDataCharacteristic  {

    constructor() {
        log.info('Mock RowerDataCharacteristic: constructor was called');
    }

    onSubscribe(maxValueSize: number, updateValueCallback: (data: Buffer) => void): void {

        log.info(`Mock RowerDataCharacteristic onSubscribe maxValueSize=${maxValueSize}`)
    }

    onUnsubscribe(): void {

        log.info(`Mock RowerDataCharacteristic onUnsubscribe`);
    }

    onData(data: any): void {

        log.info(`Mock RowerDataCharacteristic onData`);
    }
}