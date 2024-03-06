import { Data } from '../Data';
import { FitnessMachineStatusCharacteristic } from './FitnessMachineStatusCharacteristic';
import { PrimaryService } from '@abandonware/bleno';
import { RowerDataCharacteristic } from './RowerDataCharacteristic';
import { RowerFeatureCharacteristic } from './RowerFeatureCharacteristic';
import { StatusChange } from '../StatusChange';

const rowerFeatureCharacteristic = new RowerFeatureCharacteristic();
const rowerDataCharacteristic = new RowerDataCharacteristic();
const fitnessMachineStatusCharacteristic = new FitnessMachineStatusCharacteristic();

export class FitnessMachineService extends PrimaryService {
    
    rowerDataCharacteristic? : RowerDataCharacteristic | null
    rowerFeatureCharacteristic? : RowerFeatureCharacteristic | null
    fitnessMachineStatusCharacteristic?: FitnessMachineStatusCharacteristic | null;

    constructor () {
        super({
            uuid: '1826',
            characteristics: [rowerFeatureCharacteristic, rowerDataCharacteristic, fitnessMachineStatusCharacteristic]
        })
        this.rowerDataCharacteristic = rowerDataCharacteristic;
        this.rowerFeatureCharacteristic = rowerFeatureCharacteristic;
        this.fitnessMachineStatusCharacteristic = fitnessMachineStatusCharacteristic;
    }

    onData(data : Data) : void {
        this.rowerDataCharacteristic && this.rowerDataCharacteristic.onData(data);
    }

    onStatusChange(statusChange: StatusChange) {
        this.fitnessMachineStatusCharacteristic && this.fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
    }
}