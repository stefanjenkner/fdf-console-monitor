export type Data = {
    elapsedTime: number;
    distance: number;
    time500mSplit: number | null;
    time500mAverage: number | null;
    strokes: number | null;
    strokesPerMinute: number | null;
    wattsPreviousStroke: number | null;
    wattsAverage: number | null;
    caloriesPerHour: number | null;
    caloriesTotal: number | null;
    level: number;
}