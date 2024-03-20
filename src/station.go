package main

import (
	"hash/maphash"
)

var seed maphash.Seed

func init() {
	seed = maphash.MakeSeed()
}

func gomaphash(data []byte) uint64 {
	return maphash.Bytes(seed, data)
}

// parseStationLine returns an hash of the station name, the measured temperature and the station name itself
// the line is expected to be in the format "station;temperature"
func parseStationLine(line []byte) (uint64, int16, []byte) {
	// faster to begin search from the line end since the temperature should be shorter than
	// the station name
	lastIdxStationName := len(line) - 1
	for ; line[lastIdxStationName] != byteWordSeparator; lastIdxStationName-- {
	}
	stationName := line[:lastIdxStationName]
	return gomaphash(stationName), parseMeasurement(line[lastIdxStationName+1:]), stationName
}

func parseMeasurement(linePart []byte) int16 {
	measurement := int16(0)
	numberStartIdx := 0
	if linePart[0] == byteMinusSymb {
		numberStartIdx++
	}

	for i := numberStartIdx; i < len(linePart); i++ {
		if linePart[i] == byteDot {
			continue
		}
		measurement = (measurement * 10) + int16(linePart[i]-byteDigitZero)
	}

	if numberStartIdx == 1 {
		return -measurement
	}

	return measurement
}

// avgMeasurements sums the two measurements and returns the measurement value
func avgMeasurements(a, b int16) int16 {
	// return (a + b) / 2
	return (a + b) >> 1
}

// minMeasurements returns the lowest measurement value
func minMeasurements(a, b int16) int16 {
	// if a < b {
	// 	return a
	// }
	// return b
	return a + ((b - a) & ((b - a) >> 15))
}

// maxMeasurements returns the highest measurement value
func maxMeasurements(a, b int16) int16 {
	// if a > b {
	// 	return a
	// }
	// return b
	return a - ((a - b) & ((a - b) >> 15))
}
