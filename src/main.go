package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type stationLine struct {
	hash        uint64
	measurement int16
	name        []byte
}

type measurement struct {
	name        []byte
	min         int16
	max         int16
	sum         int32
	nrSightings int16
}

func printMeasurement(m measurement, w io.Writer) error {
	minHigh := m.min / 10
	minLow := m.min % 10
	if minLow < 0 {
		minLow = -minLow
	}

	avg := m.sum / int32(m.nrSightings)
	avgHigh := avg / 10
	avgLow := avg % 10
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

func readFile(inputFile *string, bufPageChan chan []byte) {
	f, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	defer close(bufPageChan)

	var buf = make([]byte, bufferSizeInBytes)

	for {
		bytesRead, err := f.Read(buf)
		if bytesRead == 0 {
			break
		}
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			panic(fmt.Errorf("reading to buffer: %w", err))
		}

		// go back to the last found new line, both in buffer as in the file
		lastNewLineIdx := lastIndexOf(buf, bytesRead, byteNewLine)

		_, err = f.Seek(int64(lastNewLineIdx-bytesRead+1), 1)
		if err != nil {
			panic(err)
		}

		// send "page" with multiple lines to line parser
		bufPageCopy := make([]byte, lastNewLineIdx)
		copy(bufPageCopy, buf[:lastNewLineIdx])
		bufPageChan <- bufPageCopy
	}
}

func parseLines(lineParseWg *sync.WaitGroup, bufPageChan chan []byte, lineChan chan stationLine) {
	defer lineParseWg.Done()
	for bufPage := range bufPageChan {
		for bufLastReadIdx := 0; bufLastReadIdx < len(bufPage); {
			lineStart, lineEnd := lineIdxs(bufLastReadIdx, bufPage, byteNewLine)
			bufLastReadIdx = lineEnd + 1

			stHash, stTemp, stName := parseStationLine(bufPage[lineStart:lineEnd])
			lineChan <- stationLine{stHash, stTemp, stName}
		}
	}
}

func shardFunc(data []byte) uint8 {
	sum := uint64(0)
	for i := 0; i < 8 && i < len(data); i++ {
		sum += uint64(data[i])
	}

	return uint8(sum % uint64(aggregatorWorkers))
}

func agregate(lineChan chan stationLine, wg *sync.WaitGroup, shardId uint8, result map[uint64]measurement) {
	defer wg.Done()
	for line := range lineChan {
		if shardFunc(line.name) == shardId {
			current := result[line.hash]
			current.name = line.name
			current.max = maxMeasurements(current.max, line.measurement)
			current.min = minMeasurements(current.min, line.measurement)
			current.sum += int32(line.measurement)
			current.nrSightings += 1
			result[line.hash] = current
		}
	}
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

	bufPageChan := make(chan []byte, 10)     // buffered channel for 100 buffer pages = 100 * bufferSize
	lineChan := make(chan stationLine, 1000) // buffered channel for 1000 station lines

	var aggregateWg sync.WaitGroup
	aggregateWg.Add(aggregatorWorkers)
	aggregatedRes := make([]map[uint64]measurement, aggregatorWorkers)
	for i := 0; i < lineParserWorkers; i++ {
		aggregatedRes[i] = make(map[uint64]measurement, nrStations/aggregatorWorkers)
		go agregate(lineChan, &aggregateWg, uint8(i), aggregatedRes[i])
	}

	var lineParseWg sync.WaitGroup
	lineParseWg.Add(lineParserWorkers)
	for i := 0; i < lineParserWorkers; i++ {
		go parseLines(&lineParseWg, bufPageChan, lineChan)
	}

	go readFile(inputFile, bufPageChan)

	lineParseWg.Wait()
	close(lineChan)

	aggregateWg.Wait()

	// agregation is done, print everything

	output := bufio.NewWriterSize(os.Stdout, bufferSizeInBytes)
	for _, res := range aggregatedRes {
		for _, v := range res {
			if len(v.name) > 0 {
				printMeasurement(v, output)
			}
		}
	}
	err := output.Flush()
	if err != nil {
		panic(err)
	}
}
