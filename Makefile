.PHONY: build generate bench

build: generate
	go build -o advent-of-code main.go

generate:
	go generate ./...

# make bench                 -> all days
# make bench BENCH=2023      -> one year
# make bench BENCH=2015/day04 -> one day
BENCH ?=
bench:
	go test ./cmd -run '^$$' -bench 'Days/$(BENCH)' -benchmem
