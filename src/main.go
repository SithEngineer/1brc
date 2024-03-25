package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

// Measuring time by:
// ❯ go build -o read .
// ❯ multitime -n 10 ./read -input measurements-1brc.txt
// multitime is available at https://tratt.net/laurie/src/multitime/

/*
            Mean        Std.Dev.    Min         Median      Max
real        6.888       0.869       5.850       6.662       8.776
user        0.005       0.001       0.002       0.005       0.006
sys         1.976       0.092       1.795       1.978       2.130
*/
// const bufferSize int = 256 * 1024 * 1024

/*
            Mean        Std.Dev.    Min         Median      Max
real        6.195       0.743       5.813       5.833       8.362
user        0.009       0.001       0.007       0.009       0.010
sys         1.889       0.074       1.769       1.893       2.053
*/
// const bufferSize int = 64 * 1024 * 1024

/*
            Mean        Std.Dev.    Min         Median      Max
real        5.864       0.086       5.796       5.827       6.106
user        0.014       0.001       0.012       0.014       0.015
sys         1.884       0.044       1.769       1.900       1.920
*/
//const bufferSize int = 32 * 1024 * 1024

/*
	Mean        Std.Dev.    Min         Median      Max

real        5.857       0.121       5.792       5.819       6.217
user        0.024       0.002       0.021       0.024       0.026
sys         1.876       0.039       1.789       1.888       1.917
*/
const bufferSize int = 16 * 1024 * 1024

/*
            Mean        Std.Dev.    Min         Median      Max
real        5.834       0.114       5.773       5.795       6.172
user        0.028       0.001       0.027       0.028       0.030
sys         1.853       0.073       1.689       1.874       1.969
*/
// const bufferSize int = 8 * 1024 * 1024

/*
            Mean        Std.Dev.    Min         Median      Max
real        6.084       0.117       6.028       6.044       6.432
user        0.052       0.001       0.050       0.052       0.054
sys         1.865       0.017       1.836       1.861       1.893
*/
// const bufferSize int = 1024 * 1024

/*
	          Mean        Std.Dev.    Min         Median      Max
real        6.168       0.103       6.108       6.128       6.469
user        0.087       0.001       0.085       0.087       0.090
sys         2.008       0.019       1.969       2.011       2.030
*/
// const bufferSize int = 512 * 1024

/*
            Mean        Std.Dev.    Min         Median      Max
real        10.086      0.203       9.923       9.956       10.563
user        2.659       0.006       2.652       2.657       2.674
sys         6.897       0.018       6.871       6.892       6.927
*/
// const bufferSize int = 1 * 1024

func main() {
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

	buf := make([]byte, bufferSize)
	for {
		bytesRead, err := f.Read(buf)
		if bytesRead == 0 {
			break
		}
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			panic(fmt.Errorf("reading to buffer: %w", err))
		}
	}
}
