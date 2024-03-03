import { RowerDataCharacteristic } from '../../fitnessmachine/RowerDataCharacteristic';

test('call updateValueCallback', () => {

    // setup
    const capture = {
        strokesPerMinute: 0,
        distance: 123,
        time500mSplit: 115,
        elapsedTime: 45,
        caloriesPerHour: 987,
        level: 0,
        watt: 105
    }
    const updateValueCallback = jest.fn();
    const rowerDataCharacteristik = new RowerDataCharacteristic()
    rowerDataCharacteristik.onSubscribe(0, updateValueCallback)
    
    // execute
    rowerDataCharacteristik.onCapture(capture);

    // verify
    const expected = Buffer.from([0x2D, 0x09, 123, 0, 0, 115, 0, 105, 0, 219, 3, 45, 0]);
    expect(updateValueCallback).toHaveBeenCalledWith(expected);
});