build:
	@go build -o ../bin/measure main.go

test:
	@go test -v
	
bench:
	# -run=^$ means run no unit test, just benchmarks
	@go test -run=^$ -bench=.

run: build
	@../bin/measure -input measurements.txt
	
time: build
	@time ../bin/measure -input measurements.txt
	
measurements-billion:
	@python3 create_measurements.py 1_000_000_000
	
measurements-million:
	@python3 create_measurements.py 1_000_000

profile-CPU: run
	@go tool pprof -http=:8080 cpu.pprof
	
profile-trace: run
	@go tool trace trace.out
