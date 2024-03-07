package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/profile"
)

type input struct {
	city string
	value float32
}

type measurements struct {
	min float32
	max float32
	mean float32
}

func parseFloat(in string) (float32) {
	f, _ := strconv.ParseFloat(in, 32)
	return float32(f)
}

func extractMeasurement(line string) (string, float32) {
	parts := strings.Split(string(line), ";")
	return parts[0], parseFloat(parts[1])
}

func printMeasurement(city string, m measurements) (string) {
	return fmt.Sprintf("%s;%.1f;%.1f;%.1f", city, m.min, m.mean, m.max)
}

func process(d input, res map[string]measurements){
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

func main()  {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop() // for CPU
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
	reader := bufio.NewReader(f)
	for{
		line, _, err := reader.ReadLine()	
		if err != nil {
			break
		}
		city, value := extractMeasurement(string(line))
		process(input{city: city, value: value}, result)
	}	
	
	for k, v := range result {
		fmt.Println(printMeasurement(k, v))
	}
}
