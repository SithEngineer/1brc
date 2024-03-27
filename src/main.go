package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Measuring time by:
// ❯ go build -o read .
// ❯ multitime -n 10 ./read -input measurements-1brc.txt
// multitime is available at https://tratt.net/laurie/src/multitime/

// reading a file
/*
            Mean        Std.Dev.    Min         Median      Max
real        5.857       0.121       5.792       5.819       6.217
user        0.024       0.002       0.021       0.024       0.026
sys         1.876       0.039       1.789       1.888       1.917
*/

// reading a file and copying the buffer into a channel
/*
            Mean        Std.Dev.    Min         Median      Max
real        7.074       1.335       5.928       6.507       9.251
user        1.156       0.005       1.149       1.156       1.164
sys         1.658       0.076       1.556       1.634       1.772
*/

const bufferSize int = 16 * 1024 * 1024

func readFile(filename string, bufPageChan chan []byte) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	defer close(bufPageChan)

	buf := make([]byte, bufferSize)

	for {
		bytesRead, err := f.Read(buf)
		if bytesRead == 0 {
			break
		}
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			panic(fmt.Errorf("reading to buffer: %w", err))
		}
		bufCopy := make([]byte, bytesRead)
		copy(bufCopy, buf[:bytesRead])
		bufPageChan <- bufCopy
	}
}

func readBufferPages(bufPageChan chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	for _ = range bufPageChan {
	}
}

func main() {
	inputFile := flag.String("input", "", "measurements file to be processed")
	flag.Parse()
	if *inputFile == "" {
		log.Fatal("measurements file not provided")
	}

	var wg sync.WaitGroup
	bufPageChan := make(chan []byte, 100)
	wg.Add(1)
	go readBufferPages(bufPageChan, &wg)
	go readFile(*inputFile, bufPageChan)
	wg.Wait()
}
