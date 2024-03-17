build:
	@go build -o ../bin/measure main.go

test:
	@go test -v --race
	
bench:
	# -run=^$ means run no unit test, just benchmarks
	@go test -bench=./... -run=^$

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

profile-CPU: run
	@go tool pprof -http=:8080 cpu.pprof

profile-mem: run
	@go tool pprof -http=:8080 mem.pprof
	
profile-trace: run
	@go tool trace trace.out
