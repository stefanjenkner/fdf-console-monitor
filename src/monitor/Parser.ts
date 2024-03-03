import { Capture } from './Capture';

export class Parser {

    public parse(message : string) : Capture {

        const totalMinutes =  Number.parseInt(message.substring(3, 5));
        const totalSeconds =  Number.parseInt(message.substring(5, 7));
        const distance = Number.parseInt(message.substring(7, 12));
        const to500mMinutes = Number.parseInt(message.substring(13, 15));
        const to500mSeconds = Number.parseInt(message.substring(15, 17));
        const strokesPerMinute = Number.parseInt(message.substring(17, 20));
        const watt = Number.parseInt(message.substring(20, 23));
        const caloriesPerHour = Number.parseInt(message.substring(23, 27));
        const level = Number.parseInt(message.substring(27, 29));

        return {
            elapsedTime: totalMinutes*60 + totalSeconds,
            time500m: to500mMinutes*60 + to500mSeconds,
            distance: distance,
            level: level,
            watts: watt,
            cals: caloriesPerHour,
            strokesPerMinute: strokesPerMinute
        }
    }
}