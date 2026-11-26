// Command benchreport turns `go test -bench Days` output from a purego run and
// a default (asm) run into a Markdown table, merging with an existing report.
// It's normally run through make:
//
//	make bench-report                  # every day
//	make bench-report BENCH=2015/day01 # just that day's row; the rest are kept
//
// or directly:
//
//	go run ./tools/benchreport purego.txt asm.txt BENCHMARKS.md
//
// Times are the median of all runs (-count=N), rounded to 2 significant
// figures. An existing row's value is only replaced when the new median differs
// by more than -threshold, so rerunning doesn't churn the git diff.
package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/perf/benchfmt"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("benchreport: ")

	threshold := flag.Float64("threshold", 0.10, "only replace an existing value if the median moved by more than this fraction")
	root := flag.String("root", ".", "repo root (used to detect days with .s files)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: benchreport [flags] purego.txt asm.txt BENCHMARKS.md")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(2)
	}
	pureFile, asmFile, outFile := flag.Arg(0), flag.Arg(1), flag.Arg(2)

	// "asm" is the default build (no -tags purego). It only runs different
	// code from purego for days with a .s file for this architecture.
	pure, err := parseBench(pureFile)
	if err != nil {
		log.Fatal(err)
	}
	asm, err := parseBench(asmFile)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := readReport(outFile)
	if err != nil {
		log.Fatal(err)
	}

	goarch := cmp.Or(asm.goarch, pure.goarch, runtime.GOARCH)
	asmDay := func(d dayID) bool { return hasAsm(*root, goarch, d) }
	merge(rows, pure, asm, *threshold, asmDay)

	report := render(rows, asm, pure, goarch, *threshold)
	if err := os.WriteFile(outFile, []byte(report), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %s (%d days)\n", outFile, len(rows))
}

// merge folds new benchmark results into rows. Days missing from both runs
// (e.g. outside a BENCH filter) are left as they are.
func merge(rows map[dayID]row, pure, asm *benchFile, threshold float64, asmDay func(dayID) bool) {
	days := map[dayID]bool{}
	for d := range pure.samples {
		days[d] = true
	}
	for d := range asm.samples {
		days[d] = true
	}

	for d := range days {
		r := rows[d]
		if s := pure.samples[d]; s != nil {
			r.purego = pick(r.purego, median(s.ns), threshold)
		}

		hasAsm := asmDay(d)
		if !hasAsm {
			r.asm = ""
		} else if s := asm.samples[d]; s != nil {
			r.asm = pick(r.asm, median(s.ns), threshold)
		}

		// allocs for what the default build actually runs
		src := pure.samples[d]
		if hasAsm && asm.samples[d] != nil {
			src = asm.samples[d]
		}
		if src != nil && len(src.allocs) > 0 {
			// allocs are whole numbers; an even -count can put the median between two
			r.allocs = strconv.Itoa(int(math.Round(median(src.allocs))))
		}
		rows[d] = r
	}
}

// pick keeps the old cell unless the new value moved by more than threshold.
func pick(old string, newNs, threshold float64) string {
	if math.IsNaN(newNs) {
		return old
	}
	if oldNs, ok := parseNs(old); ok && oldNs > 0 && math.Abs(newNs-oldNs)/oldNs <= threshold {
		return old
	}
	return formatNs(newNs)
}

// hasAsm reports whether the day's package has hand-written assembly for
// goarch (e.g. advent/2015/day1/day1_arm64.s). Assembly for other
// architectures doesn't count: the default build would just run the purego
// code again.
func hasAsm(root, goarch string, d dayID) bool {
	dir := filepath.Join(root, "advent", strconv.Itoa(d.year), fmt.Sprintf("day%d", d.day))
	matches, _ := filepath.Glob(filepath.Join(dir, "*_"+goarch+".s"))
	return len(matches) > 0
}

// ---- days ----

// dayID identifies one puzzle, e.g. 2015/day01.
type dayID struct{ year, day int }

func (d dayID) String() string { return fmt.Sprintf("%d/day%02d", d.year, d.day) }

func (d dayID) compare(o dayID) int {
	return cmp.Or(cmp.Compare(d.year, o.year), cmp.Compare(d.day, o.day))
}

// dayName matches a day, both as a report cell and as the sub-benchmark part
// of BenchmarkDays/2015/day01-16.
var dayName = regexp.MustCompile(`^(\d+)/day(\d+)$`)

// parseDay parses "2015/day01".
func parseDay(s string) (dayID, bool) {
	m := dayName.FindStringSubmatch(s)
	if m == nil {
		return dayID{}, false
	}
	// the regex only matches digits, so Atoi can't fail
	y, _ := strconv.Atoi(m[1])
	d, _ := strconv.Atoi(m[2])
	return dayID{y, d}, true
}

// ---- benchmark input ----

// sample collects one day's measurements, one entry per -count run.
type sample struct {
	ns     []float64 // time per op in nanoseconds
	allocs []float64 // allocs per op
}

// benchFile is one parsed `go test -bench` output file.
type benchFile struct {
	cpu, goarch string // from the "cpu:" and "goarch:" header lines
	samples     map[dayID]*sample
}

// maxCount is the most runs recorded for any one day (i.e. -count).
func (bf *benchFile) maxCount() int {
	n := 0
	for _, s := range bf.samples {
		n = max(n, len(s.ns))
	}
	return n
}

// parseBench reads `go test -bench Days` output. Lines that aren't benchmark
// results (PASS, ok ..., log output) are ignored by benchfmt.
func parseBench(path string) (*benchFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bf := &benchFile{samples: map[dayID]*sample{}}
	r := benchfmt.NewReader(f, path)
	for r.Scan() {
		switch rec := r.Result().(type) {
		case *benchfmt.SyntaxError:
			log.Printf("warning: %v", rec)
		case *benchfmt.Result:
			id, ok := benchDay(rec.Name)
			if !ok {
				continue
			}
			bf.cpu = cmp.Or(bf.cpu, rec.GetConfig("cpu"))
			bf.goarch = cmp.Or(bf.goarch, rec.GetConfig("goarch"))

			// benchfmt tidies ns/op into base units: sec/op.
			sec, ok := rec.Value("sec/op")
			if !ok {
				continue
			}
			s := bf.samples[id]
			if s == nil {
				s = &sample{}
				bf.samples[id] = s
			}
			s.ns = append(s.ns, sec*1e9)
			if a, ok := rec.Value("allocs/op"); ok {
				s.allocs = append(s.allocs, a)
			}
		}
	}
	return bf, r.Err()
}

// benchDay extracts the day from a Days/<year>/day<nn>[-procs] name
// (benchfmt strips the "Benchmark" prefix).
func benchDay(name benchfmt.Name) (dayID, bool) {
	base, parts := name.Parts() // "Days", ["/2015", "/day01", "-16"]
	if string(base) != "Days" || len(parts) < 2 {
		return dayID{}, false
	}
	return parseDay(strings.TrimPrefix(string(parts[0])+string(parts[1]), "/"))
}

// median returns the median of xs, or NaN if xs is empty.
func median(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	s := slices.Clone(xs)
	slices.Sort(s)
	n := len(s)
	if n%2 == 1 {
		return s[n/2]
	}
	return (s[n/2-1] + s[n/2]) / 2
}

// ---- report ----

// Report table columns. readReport and render both index cells by these.
const (
	colDay = iota
	colPurego
	colAsm
	colSpeedup
	colAllocs
	numCols
)

// header is the report table's header row.
var header = [numCols]string{
	colDay:     "Day",
	colPurego:  "purego",
	colAsm:     "asm",
	colSpeedup: "speedup",
	colAllocs:  "allocs",
}

// row holds a report row's cells as rendered text, not numbers, so a value
// kept by pick is written back byte for byte and doesn't show up in the diff.
// speedup isn't stored; render derives it from purego and asm.
type row struct {
	purego, asm, allocs string
}

// readReport loads rows from an existing report's table, if any.
func readReport(path string) (map[dayID]row, error) {
	rows := map[dayID]row{}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return rows, nil
	}
	if err != nil {
		return nil, err
	}
	for line := range strings.Lines(string(data)) {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "|") {
			continue
		}
		cells := strings.Split(strings.Trim(line, "|"), "|")
		if len(cells) != numCols {
			continue
		}
		for i := range cells {
			cells[i] = strings.TrimSpace(cells[i])
		}
		d, ok := parseDay(cells[colDay])
		if !ok {
			continue // header / separator
		}
		rows[d] = row{
			purego: cells[colPurego],
			asm:    cells[colAsm],
			allocs: cells[colAllocs],
		}
	}
	return rows, nil
}

// render writes the whole report: a short header describing the machine and
// method, then the table sorted by day.
func render(rows map[dayID]row, asm, pure *benchFile, goarch string, threshold float64) string {
	cpu := cmp.Or(asm.cpu, pure.cpu, "unknown cpu")
	count := max(asm.maxCount(), pure.maxCount())

	var b strings.Builder
	b.WriteString("# Benchmarks\n\n")
	// runtime.Version is the Go that ran this tool. Under `make bench-report`
	// that's the same toolchain `go test` used for the benchmarks.
	fmt.Fprintf(&b, "%s · %s · %s · median of %d runs, 2 significant figures\n\n", cpu, goarch, runtime.Version(), count)
	fmt.Fprintf(&b, "Generated by `make bench-report`. *asm* is only shown for days with hand-written assembly. "+
		"A value only changes when its median moves more than %.0f%%.\n\n", threshold*100)

	table := [][numCols]string{header}
	for _, d := range slices.SortedFunc(maps.Keys(rows), dayID.compare) {
		r := rows[d]
		var speedup string
		if p, ok := parseNs(r.purego); ok {
			if a, ok := parseNs(r.asm); ok && a > 0 {
				speedup = strconv.FormatFloat(roundSig(p/a, 2), 'f', -1, 64) + "×"
			}
		}
		table = append(table, [numCols]string{
			colDay:     d.String(),
			colPurego:  r.purego,
			colAsm:     r.asm,
			colSpeedup: speedup,
			colAllocs:  r.allocs,
		})
	}

	var widths [numCols]int
	for _, cells := range table {
		for i, c := range cells {
			widths[i] = max(widths[i], utf8.RuneCountInString(c))
		}
	}

	// fmt pads %*s by runes, so µ and × line up. The Day column is left
	// aligned, the numbers right aligned.
	writeRow := func(cells [numCols]string) {
		b.WriteString("|")
		for i, c := range cells {
			if i == colDay {
				fmt.Fprintf(&b, " %-*s |", widths[i], c)
			} else {
				fmt.Fprintf(&b, " %*s |", widths[i], c)
			}
		}
		b.WriteString("\n")
	}
	writeRow(table[0])
	b.WriteString("|")
	for i, w := range widths {
		sep := strings.Repeat("-", w)
		if i != colDay {
			sep = sep[:w-1] + ":" // right-align numeric columns
		}
		fmt.Fprintf(&b, " %s |", sep)
	}
	b.WriteString("\n")
	for _, cells := range table[1:] {
		writeRow(cells)
	}
	return b.String()
}

// ---- number formatting ----

// roundSig rounds v to n significant figures. It scales by an exact integer
// power of ten so results like 5300 come out exact (dividing by 0.01 doesn't).
func roundSig(v float64, n int) float64 {
	if v == 0 {
		return 0
	}
	k := n - int(math.Ceil(math.Log10(math.Abs(v))))
	if k >= 0 {
		p := math.Pow10(k)
		return math.Round(v*p) / p
	}
	p := math.Pow10(-k)
	return math.Round(v/p) * p
}

// formatNs renders a duration in ns with 2 significant figures: 720ns, 5.3µs, 120µs, 1.5s.
func formatNs(ns float64) string {
	return time.Duration(roundSig(ns, 2)).String()
}

// parseNs is the inverse of formatNs. ok is false for blank/unknown cells.
func parseNs(s string) (float64, bool) {
	d, err := time.ParseDuration(strings.TrimSpace(s))
	if err != nil {
		return 0, false
	}
	return float64(d), true
}
