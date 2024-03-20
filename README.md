# The one billion row challenge

[Based on the original](https://github.com/gunnarmorling/1brc/blob/main/README.md#1%EF%B8%8F%E2%83%A3%EF%B8%8F-the-one-billion-row-challenge)

## Appoach #1: Make it work!

There are a few objectives for a first approach:
1. Read the filename from the command args;
2. Open the file and read it, line by line;
3. Make the code simple and readable. This also means testable, so better by starting with tests;
4. Don't allow the presure to make it fast corrupt the readability of the code;

1. and 2. are easy but then it starts to become hard to do 3. when the focus is "how fast can I make this?!" so that is why 4. is an objective that must be fufilled, for the first approach.

Check the [rules](https://github.com/gunnarmorling/1brc/blob/main/README.md) to see how to handle that file input.

> * Input value ranges are as follows:
>   * Station name: non null UTF-8 string of min length 1 character and max length 100 bytes, containing neither `;` nor `\n` characters. (i.e. this could be 100 one-byte characters, or 50 two-byte characters, etc.)
>   * Temperature value: non null double between -99.9 (inclusive) and 99.9 (inclusive), always with one fractional digit
> 
> * There is a maximum of 10,000 unique station names
>
> * Line endings in the file are `\n` characters on all platforms

Build the code to successfully process the data... And time it! `time ...`.

Now that we established a baseline in terms of code and time needed, lets begin our journey.

## How to profile a Go program

1. Write benchmark tests and `go test -bench=.`
2. Use Go pprof by profiling the code, then (for a CPU profile) `go tool pprof -http=:8080 cpu.pprof` and open the browser
3. Similar to 2. but create a Trace profile, then run `go tool trace trace.out` and open a web browser

## Optimization journey starts here

### Baseline

From the first execution, which took 156.67 seconds. A CPU profiling makes it is clear that reading a file line by line is **very** time consuming. More than 70% of the CPU is wasted in syscalls, from reading a file, line by line. Interestingly we have a similar issue with writing to the standard output.

### Using buffers

Knowing that each station name has a max of 100 bytes + `;` + reading (longest possible is `-99.9`) = 106 bytes for a maximum of each line. This means that we can do three things:
1. Read a chunk of the file into a buffer
2. Find the last new line to seek back in the file reading
3. Read line by line using a smaller buffer

All this work to reduce time to 149.62 seconds using a 1MB buffer. Sad. Looking at the computation profile a lot of CPU time is invested in `extractMeasurement` and `process` functions, around 20% and 25% respectively, so this gives us two targets for the next iteration. There is one detail though: `process` seems to spend most of it's time - 20% - on a map access, let's see how we can change that.

### Using smaller storage units for data

`-99.9` can be decomposed in 13 bits:
1. `-`  = 1 bit
2. `99` = 7 bits
3. `.`  = 1 bit
4. `9`  = 4 bit

Using an uint16 with big endian we get `-99.9` = 1001 1111 1001 1001
Using an uint16 with little endian we get `-99.9` = 1001 1001 1111 1001

13 bits is the maximum a single temperature measurement will use, so we can use two `byte`s or a single `uint16` but this also means (dangerous) bitwise operations.

However if we use a number from 0 to 999, thus droping the decimal separator and keeping the signal bit:
[signal][0..999] = 11 bits thus `-99.9` becomes 1 1111100111

In go we would use a uint16 to store the 11 bits nevertheles, so we can take another shortcut for the signal.
if 'uint16' stores values between '0 - 1024' then 'int16' stores values between '-512 - 511'. This means our avg, min and max operations will use int16 instead of float32, and to print the value we can just do 'number / 10' for the integer part and the floating point is positioned by 'number % 10'.

Parsing a measurement is done by reading byte by byte and aggregating that to a result using the formula: existing result = (read byte - byte '0') + (existing result * 10). If the first byte is a '-' the result is a negative value, otherwise the result is a positive number.

The average between two measurements consists in two operations:
1. Adding the existing and new measurement
2. Divide the previous result by 2

Point 2. can be achieved by shifting bits to the right once.

Cities are strings and a string in go is an imutable array of bytes. Let's keep them like this at least for now.

At this point the time to process the whole file is around 63 seconds.

### Challenges to solve

#### Map access & assignment

There was a margninal gain by using an hashed station name as a uint64 instead of a string for a map key. This forced to store the station name in the value entry, making the value bigger than necessary.

#### `lineIdxs`

Used to calculate the beginning and end of a line, in a buffer "page". Consumes ~22% of CPU cycles.

#### `lineIdxs`

Used to calculate the beginning and end of a line, in a buffer "page". Consumes ~22% of CPU cycles. 

#### `parseStationLine` 

After knowing the beginning and end of a line, this function is used to parse the line data. Consumers ~20% of CPU cycles.

#### `mapaccess1_fast64` and `mapassign_fast64`

~31% and ~7% respectively CPU cycles consumed. This can be mitigated by not using a map... Or not using go map.
