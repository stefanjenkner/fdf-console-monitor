import { Parser } from "../monitor/Parser";

test("parsing raw data to captur object", () => {

    // setup
    const parser = new Parser();

    // execute
    const capture = parser.parse("A8000060001410243028105065904")

    // verify
    expect(capture.elapsedTime).toBe(6);
    expect(capture.distance).toBe(14);
    expect(capture.watt).toBe(105);
    expect(capture.caloriesPerHour).toBe(659);
    expect(capture.level).toBe(4);
});