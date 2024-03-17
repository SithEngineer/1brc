package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const newLine byte = 10
const bufferSize int = 1 * 1024 * 1024

type input struct {
	city  string
	value float32
}

type measurements struct {
	min  float32
	max  float32
	mean float32
}

func extractMeasurement(line []byte) (string, float32, error) {
	parts := strings.Split(string(line), ";")
	f, err := strconv.ParseFloat(parts[1], 32)
	if err != nil {
		return "", 0.0, err
	}
	return parts[0], float32(f), nil
}

func printMeasurement(city string, m measurements, w io.Writer) error {
	_, err := fmt.Fprintf(w, "%s;%.1f;%.1f;%.1f\n", city, m.min, m.mean, m.max)
	return err
}

func process(d input, res map[string]measurements) {
	entry, ok := res[d.city]
	if !ok {
		entry = measurements{min: d.value, max: d.value, mean: d.value}
	} else {
		entry.mean += d.value
		entry.mean /= 2
		if d.value > entry.max {
			entry.max = d.value
		}
		if d.value < entry.min {
			entry.min = d.value
		}
	}
	res[d.city] = entry
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

	result := map[string]measurements{}

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

			city, value, err := extractMeasurement(buf[lineStart:lineEnd])
			if err != nil {
				panic(fmt.Errorf("reading line: %w", err))
			}
			process(input{city: city, value: value}, result)
		}
	}

	output := bufio.NewWriter(os.Stdout)
	for k, v := range result {
		printMeasurement(k, v, output)
	}
	err = output.Flush()
	if err != nil {
		panic(err)
	}
}
