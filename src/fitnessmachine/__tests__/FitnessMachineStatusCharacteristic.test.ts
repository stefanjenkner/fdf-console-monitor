import { FitnessMachineStatusCharacteristic } from '../FitnessMachineStatusCharacteristic';
import { StatusChange } from '../../StatusChange';

describe('FitnessMachineStatusCharacteristic', () => {
    const updateValueCallback = jest.fn();
    const fitnessMachineStatusCharacteristic = new FitnessMachineStatusCharacteristic();

    beforeEach(() => {
        fitnessMachineStatusCharacteristic.onSubscribe(0, updateValueCallback);
        updateValueCallback.mockReset();
    });

    test('should call updateValueCallback on Started', () => {
        // setup
        const statusChange = StatusChange.Started;
        // execute
        fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
        // verify
        const expected = Buffer.from([0x04, 0x00]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });
    
    test('should call updateValueCallback on PausedOrStopped', () => {
        // setup
        const statusChange = StatusChange.PausedOrStopped;
        // execute
        fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
        // verify
        const expected = Buffer.from([0x02, 0x00]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });

    test('should call updateValueCallback on Reset', () => {
        // setup
        const statusChange = StatusChange.Reset;
        // execute
        fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
        // verify
        const expected = Buffer.from([0x01, 0x00]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });

    test('should not call updateValueCallback on unsupported operation', () => {
        // setup
        const statusChange = StatusChange.LevelChanged;
        // execute
        fitnessMachineStatusCharacteristic.onStatusChange(statusChange);
        // verify
        expect(updateValueCallback).not.toHaveBeenCalled();
    });
});