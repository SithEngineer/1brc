build:
	@go build -o ../bin/measure .

test:
	@go test -v --race
	
bench:
	@go test -bench=.

run: build
	@../bin/measure -input measurements.txt
	
time: build
	@time ../bin/measure -input measurements.txt

time-1brc: build
	@time ../bin/measure -input measurements-1brc.txt
	
measurements-billion:
	@python3 create_measurements.py 1_000_000_000
	
measurements-million:
	@python3 create_measurements.py 1_000_000

profile-cpu: run
	@go tool pprof -http=:8080 cpu.pprof

profile-mem: run
	@go tool pprof -http=:8080 mem.pprof
	
profile-trace: run
	@go tool trace trace.out
