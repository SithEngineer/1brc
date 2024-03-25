build:
	@go build -o ../bin/measure .
	
run: build
	@../bin/measure -input measurements.txt
	
time: build
	@multitime -n 10 ../bin/measure -input measurements.txt

time-1brc: build
	@multitime -n 10  ../bin/measure -input measurements-1brc.txt
	
measurements-billion:
	@python3 create_measurements.py 1_000_000_000
	
measurements-million:
	@python3 create_measurements.py 1_000_000
