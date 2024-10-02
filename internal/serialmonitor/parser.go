package serialmonitor

import "strconv"

func parse(line string) capture {
	totalMinutes, _ := strconv.ParseUint(line[3:5], 10, 64)
	totalSeconds, _ := strconv.ParseUint(line[5:7], 10, 64)
	distance, _ := strconv.ParseUint(line[7:12], 10, 64)
	to500mMinutes, _ := strconv.ParseUint(line[13:15], 10, 64)
	to500mSeconds, _ := strconv.ParseUint(line[15:17], 10, 64)
	strokesPerMinute, _ := strconv.ParseUint(line[17:20], 10, 64)
	watt, _ := strconv.ParseUint(line[20:23], 10, 64)
	caloriesPerHour, _ := strconv.ParseUint(line[23:27], 10, 64)
	level, _ := strconv.ParseUint(line[27:29], 10, 64)

	return capture{
		elapsedTime:      totalMinutes*60 + totalSeconds,
		distance:         distance,
		time500m:         to500mMinutes*60 + to500mSeconds,
		strokesPerMinute: strokesPerMinute,
		watts:            watt,
		cals:             caloriesPerHour,
		level:            level,
	}
}
