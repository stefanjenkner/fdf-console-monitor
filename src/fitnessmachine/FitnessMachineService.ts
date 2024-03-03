import { PrimaryService } from '@abandonware/bleno';
import { RowerFeatureCharacteristic } from './RowerFeatureCharacteristic';
import { RowerDataCharacteristic } from './RowerDataCharacteristic';
import { Data } from '../monitor/Data';

const rowerFeatureCharacteristic = new RowerFeatureCharacteristic();
const rowerDataCharacteristic = new RowerDataCharacteristic();

export class FitnessMachineService extends PrimaryService {

    rowerDataCharacteristic? : RowerDataCharacteristic | null
    rowerFeatureCharacteristic? : RowerFeatureCharacteristic | null

    constructor () {
        super({
            uuid: '1826',
            characteristics: [rowerFeatureCharacteristic, rowerDataCharacteristic]
        })
        this.rowerDataCharacteristic = rowerDataCharacteristic;
        this.rowerFeatureCharacteristic = rowerFeatureCharacteristic;
    }

    onData(data : Data) : void {

        this.rowerDataCharacteristic && this.rowerDataCharacteristic.onData(data);
    }
}