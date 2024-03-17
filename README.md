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

...all this to reduce time to 149.62 seconds using a 1MB buffer... 
