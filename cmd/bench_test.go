package cmd

import (
	"fmt"
	"os"
	"testing"
)

// BenchmarkDays benchmarks every registered day against its real input in
// inputs/<year>/day<n>.txt. Days without a downloaded input are skipped.
//
//	go test ./cmd -run '^$' -bench 'Days'            # all days
//	go test ./cmd -run '^$' -bench 'Days/2023'       # one year
//	go test ./cmd -run '^$' -bench 'Days/2015/day04' # one day
func BenchmarkDays(b *testing.B) {
	for _, r := range getAllRunners() {
		b.Run(fmt.Sprintf("%d/day%02d", r.year, r.day), func(b *testing.B) {
			// tests run with the package dir (cmd/) as the working directory
			path := fmt.Sprintf("../inputs/%d/day%d.txt", r.year, r.day)
			input, err := os.ReadFile(path)
			if err != nil {
				b.Skipf("no input: %v", err)
			}

			// Some days mutate their input in place (e.g. 2015 day 11), so each
			// iteration gets a fresh copy. The copy is a single memcpy and is
			// cheaper and less noisy than StopTimer/StartTimer per iteration.
			buf := make([]byte, len(input))
			b.ReportAllocs()
			for b.Loop() {
				copy(buf, input)
				if _, _, err := r.Run(buf); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
