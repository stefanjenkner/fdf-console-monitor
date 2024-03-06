import { Data } from '../Data';
import { FitnessMachineStatusCharacteristic } from './FitnessMachineStatusCharacteristic';
import { PrimaryService } from '@abandonware/bleno';
import { RowerDataCharacteristic } from './RowerDataCharacteristic';
import { RowerFeatureCharacteristic } from './RowerFeatureCharacteristic';
import { StatusChange } from '../StatusChange';

export class FitnessMachineService extends PrimaryService {
    
    rowerDataCharacteristic : RowerDataCharacteristic;
    fitnessMachineStatusCharacteristic: FitnessMachineStatusCharacteristic;

    constructor () {
        const rowerDataCharacteristic = new RowerDataCharacteristic();
        const rowerFeatureCharacteristic = new RowerFeatureCharacteristic();
        const fitnessMachineStatusCharacteristic = new FitnessMachineStatusCharacteristic();
        super({
            uuid: '1826',
            characteristics: [rowerFeatureCharacteristic, rowerDataCharacteristic, fitnessMachineStatusCharacteristic]
        })
        this.rowerDataCharacteristic = rowerDataCharacteristic;
        this.fitnessMachineStatusCharacteristic = fitnessMachineStatusCharacteristic;
    }

    onData(data : Data) : void {
        this.rowerDataCharacteristic && this.rowerDataCharacteristic.onData(data);
    }

    onStatusChange(statusChange: StatusChange) {
        this.fitnessMachineStatusCharacteristic && this.fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
    }
}