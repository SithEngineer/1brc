package main

const bufferSizeInBytes int = 4_000_000
const lineParserWorkers = 3
const aggregatorWorkers = 4

const nrStations int = 10000

const byteNewLine byte = 10
const byteWordSeparator byte = ';'
const byteMinusSymb byte = '-'
const byteDot byte = '.'
const byteDigitZero byte = '0'
