import { Data } from '../../Data';
import { RowerDataCharacteristic } from '../RowerDataCharacteristic';

describe('call updateValueCallback', () => {
    const updateValueCallback = jest.fn();
    const rowerDataCharacteristik = new RowerDataCharacteristic();

    beforeEach(() => {
        rowerDataCharacteristik.onSubscribe(0, updateValueCallback);
    });

    test('when active', () => {

        // setup
        const data: Data = {
            strokes: 23,
            strokesPerMinute: 31,
            distance: 123,
            remainingDistance: null,
            time500mSplit: 115,
            time500mAverage: null,
            elapsedTime: 45,
            caloriesPerHour: 987,
            caloriesTotal: null,
            wattsPreviousStroke: 105,
            wattsAverage: null,
            level: 0,
        };
        rowerDataCharacteristik.onSubscribe(0, updateValueCallback);

        // execute
        rowerDataCharacteristik.onData(data);

        // verify
        const expected = Buffer.from([0x2C, 0x08, 62, 23, 0, 123, 0, 0, 115, 0, 105, 0, 45, 0]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });

    test('when paused or stopped', () => {

        // setup
        const data: Data = {
            strokes: 45,
            strokesPerMinute: 0,
            distance: null,
            remainingDistance: 877,
            time500mSplit: null,
            time500mAverage: null,
            elapsedTime: 45,
            caloriesPerHour: null,
            caloriesTotal: null,
            wattsPreviousStroke: null,
            wattsAverage: null,
            level: 0,
        };

        // execute
        rowerDataCharacteristik.onData(data);

        // verify
        const expected = Buffer.from([0x01, 0x08, 45, 0]);
        expect(updateValueCallback).toHaveBeenCalledWith(expected);
    });
});