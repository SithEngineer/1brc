package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type measurement struct {
	name []byte
	min  int16
	avg  int16
	max  int16
}

func printMeasurement(m measurement, w io.Writer) error {
	minHigh := m.min / 10
	minLow := m.min % 10
	if minLow < 0 {
		minLow = -minLow
	}

	avgHigh := m.avg / 10
	avgLow := m.avg % 10
	if avgLow < 0 {
		avgLow = -avgLow
	}

	maxHigh := m.max / 10
	maxLow := m.max % 10
	if maxLow < 0 {
		maxLow = -maxLow
	}

	_, err := fmt.Fprintf(w, "%s;%d.%d;%d.%d;%d.%d\n", m.name, minHigh, minLow, avgHigh, avgLow, maxHigh, maxLow)
	if err != nil {
		return err
	}
	return nil
}

func lastIndexOf(b []byte, maxIdx int, delym byte) int {
	for i := maxIdx - 1; i >= 0; i-- {
		if b[i] == delym {
			return i
		}
	}
	return -1
}

func lineIdxs(sourceStartIdx int, source []byte, delym byte) (int, int) {
	start := sourceStartIdx
	end := sourceStartIdx
	for ; end < len(source); end++ {
		if source[end] == delym {
			break
		}
	}
	return start, end
}

func main() {
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop() // for CPU
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop() // for memory
	// defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop() // for trace

	inputFile := flag.String("input", "", "measurements file to be processed")
	flag.Parse()
	if *inputFile == "" {
		log.Fatal("measurements file not provided")
	}

	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	result := make(map[uint64]measurement, 10000)

	buf := make([]byte, bufferSize)
	for {
		bytesRead, err := f.Read(buf)
		if bytesRead == 0 {
			break
		}
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			panic(fmt.Errorf("reading to buffer: %w", err))
		}

		// go back to the last found new line, both in buffer as in the file
		lastNewLineIdx := lastIndexOf(buf, bytesRead, newLine)

		if bytesRead == bufferSize {
			_, err := f.Seek(int64(lastNewLineIdx-bytesRead+1), 1)
			if err != nil {
				panic(err)
			}
		}

		for bufLastReadIdx := 0; bufLastReadIdx < lastNewLineIdx; {
			lineStart, lineEnd := lineIdxs(bufLastReadIdx, buf, newLine)
			bufLastReadIdx = lineEnd + 1

			stHash, stTemp, stName := parseStationLine(buf[lineStart:lineEnd])
			current := result[stHash]
			current.name = stName
			current.max = maxMeasurements(current.max, stTemp)
			current.min = minMeasurements(current.min, stTemp)
			current.avg = avgMeasurements(current.avg, stTemp)
			result[stHash] = current
		}
	}

	output := bufio.NewWriterSize(os.Stdout, bufferSize)
	for _, v := range result {
		if len(v.name) > 0 {
			printMeasurement(v, output)
		}
	}
	err = output.Flush()
	if err != nil {
		panic(err)
	}
}
