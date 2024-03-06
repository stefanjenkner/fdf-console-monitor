import { Data } from './Data';
import { StatusChange } from './StatusChange';

export type MonitorEvents = {
    'connect': (err: Error | null) => void;
    'disconnect': (err: Error | null) => void;
    'data': (data: Data) => void;
    'statusChanged': (statusChange: StatusChange) => void;
};
