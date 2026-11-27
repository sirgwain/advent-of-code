.PHONY: build build-purego generate bench bench-report

build: generate
	go build -o advent-of-code main.go

build-purego:
	go build -tags purego -o advent-of-code main.go

generate:
	go generate ./...

# make bench                 -> all days
# make bench BENCH=2023      -> one year
# make bench BENCH=2015/day04 -> one day
BENCH ?=
bench:
	go test ./cmd -run '^$$' -bench 'Days/$(BENCH)' -benchmem

# make bench-report                  -> rerun all days, update BENCHMARKS.md
# make bench-report BENCH=2015/day01 -> only update that day's row
# Each day runs COUNT times in both the purego and default (asm) builds.
COUNT ?= 5
bench-report:
	@mkdir -p .bench
	go test ./cmd -run '^$$' -bench 'Days/$(BENCH)' -benchmem -count=$(COUNT) -tags purego | tee .bench/purego.txt
	go test ./cmd -run '^$$' -bench 'Days/$(BENCH)' -benchmem -count=$(COUNT) | tee .bench/asm.txt
	go run ./tools/benchreport .bench/purego.txt .bench/asm.txt BENCHMARKS.md
