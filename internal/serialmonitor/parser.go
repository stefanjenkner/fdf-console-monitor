package serialmonitor

import (
	"errors"
	"math"
	"strconv"
)

func parse(line string) capture {
	totalMinutes, _ := strconv.ParseUint(line[3:5], 10, 8)
	totalSeconds, _ := strconv.ParseUint(line[5:7], 10, 8)
	elapsedTime, _ := getSecondsUint16(totalMinutes, totalSeconds)
	distance, _ := strconv.ParseUint(line[7:12], 10, 16)
	to500mMinutes, _ := strconv.ParseUint(line[13:15], 10, 8)
	to500mSeconds, _ := strconv.ParseUint(line[15:17], 10, 8)
	time500m, _ := getSecondsUint16(to500mMinutes, to500mSeconds)
	strokesPerMinute, _ := strconv.ParseUint(line[17:20], 10, 8)
	watt, _ := strconv.ParseUint(line[20:23], 10, 16)
	caloriesPerHour, _ := strconv.ParseUint(line[23:27], 10, 16)
	level, _ := strconv.ParseUint(line[27:29], 10, 8)

	return capture{
		elapsedTime:      elapsedTime,
		distance:         uint16(distance),
		time500m:         time500m,
		strokesPerMinute: uint8(strokesPerMinute),
		watts:            uint16(watt),
		cals:             uint16(caloriesPerHour),
		level:            uint8(level),
	}
}

var (
	SecondsOutOfRange = errors.New("seconds out of range")
)

func getSecondsUint16(minutes uint64, seconds uint64) (uint16, error) {

	result := minutes*60 + seconds
	if result > math.MaxUint16 {
		return 0, SecondsOutOfRange
	}

	return uint16(result), nil
}
