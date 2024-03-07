import { FitnessMachineStatusCharacteristic } from '../FitnessMachineStatusCharacteristic';
import { StatusChange } from '../../StatusChange';

describe('call updateValueCallback', () => {
    const updateValueCallback = jest.fn();
    const fitnessMachineStatusCharacteristic = new FitnessMachineStatusCharacteristic();

    beforeEach(() => {
        fitnessMachineStatusCharacteristic.onSubscribe(0, updateValueCallback);
    });

    test('when Started', () => {
        // setup
        const statusChange = StatusChange.Started;
        // execute
        fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
        // verify
        const expected = Buffer.from([0x04, 0x00]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });
    
    test('when PausedOrStopped', () => {
        // setup
        const statusChange = StatusChange.PausedOrStopped;
        // execute
        fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
        // verify
        const expected = Buffer.from([0x02, 0x00]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });

    test('when Reset', () => {
        // setup
        const statusChange = StatusChange.Reset;
        // execute
        fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
        // verify
        const expected = Buffer.from([0x01, 0x00]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });
});