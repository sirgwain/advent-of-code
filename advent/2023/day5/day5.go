package day5

import (
	"bytes"
	"fmt"
	"math"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code/advent"
	"github.com/sirgwain/advent-of-code/advent/color"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var seedRangeStyle = lipgloss.NewStyle().Foreground(color.BrightGreen82)
var seedRangeMappingSourceStyle = lipgloss.NewStyle().Foreground(color.VeryLightMalachiteGreen84)
var mappingMatchedStyle = lipgloss.NewStyle().Foreground(color.PureMagenta200)
var mappingUnmatchedStyle = lipgloss.NewStyle().Foreground(color.BrightRed124)

var p = message.NewPrinter(language.English)

type Day struct {
}

type mapping struct {
	name   string
	ranges []rng
}

type rng struct {
	dst  int
	src  int
	size int
}

func (r rng) String() string {
	if r.dst == 0 {
		return p.Sprintf("rng: %d-%d (%d)", r.src, r.src+r.size-1, r.size)
	}
	return p.Sprintf("rng: src: %d, dst: %d, size: %d", r.src, r.dst, r.size)
}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	seedsBytes, mappingBytes, _ := bytes.Cut(bytes.TrimSpace(input), []byte{'\n'})
	seedsBytes = seedsBytes[len("seeds: "):]
	seeds := make([]int, 0, 20)
	for s := range bytes.FieldsSeq(seedsBytes) {
		// skip the space
		num, _ := advent.Atoi(s)
		seeds = append(seeds, num)
	}

	// record the seed ranges
	seedRanges := make([]rng, 0, 20)
	for i := 0; i < len(seeds); i += 2 {
		seedRanges = append(seedRanges, rng{src: seeds[i], size: seeds[i+1]})
	}

	// a mapping looks like this
	// seed-to-soil map:
	// 50 98 2
	// 52 50 48
	//
	// where the numbers are dst src range
	mappings := []mapping{}
	for m := range bytes.SplitSeq(mappingBytes, []byte("\n\n")) {

		// skip the first line, it's just a descriptor
		// every remaining line is 3 numbers
		name, lines, _ := bytes.Cut(bytes.TrimSpace(m), []byte{'\n'})
		var ranges []rng
		for line := range bytes.SplitSeq(lines, []byte{'\n'}) {
			dst, numLen := advent.Atoi(line)
			line = line[numLen+1:]
			src, numLen := advent.Atoi(line)
			line = line[numLen+1:]
			size, numLen := advent.Atoi(line)

			ranges = append(ranges, rng{dst: dst, src: src, size: size})
		}

		mappings = append(mappings, mapping{name: string(name), ranges: ranges})
	}

	// part one, single seed mapped to ranges
	s1 = math.MaxInt
	for i := range seeds {
		sMapped := 0
		for _, m := range mappings {
			for _, r := range m.ranges {
				sMapped = mapSeedToRange(seeds[i], r)
				if seeds[i] != sMapped {
					// record the new location and break
					seeds[i] = sMapped
					break
				}
			}
		}

		// min of the last mapping (location) is the solution
		s1 = min(s1, sMapped)
	}

	// part 2, map seed ranges to mapping ranges
	s2 = math.MaxInt
	for _, sr := range seedRanges {

		output := []rng{sr}
		for _, m := range mappings {
			if advent.Viz.Enabled {
				fmt.Println(seedRangeMappingSourceStyle.Render(m.name))
				fmt.Println(seedRangeStyle.Render("SEED RANGE:"))
				for _, osr := range output {
					fmt.Println(seedRangeStyle.Render(seedRangeString(osr.src, osr.size)))
				}
			}
			newOutput := []rng{}
			for _, osr := range output {
				processed, ok := processSeedRange(osr, m.ranges)
				newOutput = append(newOutput, processed...)
				if ok && advent.Viz.Enabled {
					fmt.Println(seedRangeStyle.Render("OUTPUT:"))
					for _, osr := range newOutput {
						fmt.Println(seedRangeStyle.Render(seedRangeString(osr.src, osr.size)))
					}
				}
			}
			output = newOutput
		}

		for _, osr := range output {
			s2 = min(s2, osr.src)
		}

		if advent.Viz.Enabled {
			for _, sr := range output {
				fmt.Println(seedRangeStyle.Render("AFTER MAPPING:"))
				fmt.Println(seedRangeStyle.Render(seedRangeString(sr.src, sr.size)))
			}
		}
	}

	return s1, s2, nil
}

// processSeedRange takes a single seed range and maps it to a mapping
// returns the mapped seed ranges
func processSeedRange(s rng, ranges []rng) ([]rng, bool) {
	output := []rng{}
	mapped := false
	for i, r := range ranges {
		if sMapped, leftover, ok := mapSeedRangeToRange(s, r); ok {
			if advent.Viz.Enabled && ok {
				fmt.Println(mappingMatchedStyle.Render(seedRangeString(r.src, r.size)))
				fmt.Println(seedRangeMappingSourceStyle.Render(seedRangeString(r.dst, r.size)))
				fmt.Println(seedRangeStyle.Render(seedRangeString(sMapped.src, sMapped.size)))

			}
			output = append(output, sMapped)
			if len(leftover) > 0 {
				if advent.Viz.Enabled {
					fmt.Println(seedRangeStyle.Render("LEFTOVER:"))
				}
			}
			// send the leftover through the remaning ranges in this mapping
			for _, lsr := range leftover {
				if advent.Viz.Enabled {
					fmt.Println(seedRangeStyle.Render(seedRangeString(lsr.src, lsr.size)))
				}
				subOutput, _ := processSeedRange(lsr, ranges[i+1:])
				output = append(output, subOutput...)
			}
			mapped = true
			// after mapping, we are done because the recursive call
			// to leftovers handles mapping the remaining ranges
			break
		}

		if !mapped && advent.Viz.Enabled {
			fmt.Println(mappingUnmatchedStyle.Render(seedRangeString(r.src, r.size)))
		}
	}

	if !mapped {
		// per the puzzle, any seed not in a range stays in place
		return []rng{s}, false
	}

	return output, mapped
}

// mapSeedToRange maps a seed to a range, or returns the seed if it doesn't fall within the range
func mapSeedToRange(seed int, r rng) int {
	// range is dst: 50, src: 98, size: 2
	// so 51 is between 50 - 52
	if seed >= r.src && seed < r.src+r.size {
		return r.dst + seed - r.src
	}
	return seed
}

// mapSeedRangeToRange maps a range of seends onto a range mapping
// returns first the mapping that falls within the range, then leftover parts, if any
func mapSeedRangeToRange(seed rng, r rng) (rng, []rng, bool) {
	sStart := seed.src
	sEnd := seed.src + seed.size - 1

	rStart := r.src
	rEnd := r.src + r.size - 1

	if sStart < rStart && sEnd > rEnd {
		// our seed completely overlaps a region, so split it into three parts
		// 0123456789
		// |--------|	seed: 0-9 (10)
		//    |--|		src: 3-6 (4)
		return rng{src: r.dst, size: r.size},
			[]rng{
				{src: sStart, size: rStart - sStart},
				{src: rEnd + 1, size: sEnd - rEnd},
			}, true
	}

	if sStart >= rStart && sStart <= rEnd {
		// this seed range falls somewhere within the range after the start
		overflow := sEnd - rEnd
		if overflow <= 0 {
			// we have no overflow
			// 0123456789
			//   |----|		seed: 2-7 (6)
			// |--------|	src: 0-9 (10)
			return rng{src: r.dst + (seed.src - r.src), size: seed.size}, nil, true
		}

		// we have overflow like this
		// 0123456789
		//   |----|		seed: 2-7 (6)
		// |---|		src: 0-4 (5)
		overlapStart := sStart - rStart
		return rng{src: r.dst + (overlapStart), size: r.size - overlapStart}, []rng{{src: rEnd + 1, size: overflow}}, true
	}

	if sEnd >= rStart && sEnd <= rEnd {
		// this seed falls within the range before the start
		// 0123456789
		// |---|		seed: 0-4 (5)
		//   |----|		src: 2-7 (6)
		overlap := sEnd - rStart + 1
		return rng{src: r.dst, size: overlap}, []rng{{src: seed.src, size: seed.size - overlap}}, true
	}

	return seed, nil, false
}

func seedRangeString(start, size int) string {
	if size > 100 {
		// can't print huge sizes, just summarize
		return p.Sprintf("%d - %d (%d)", start, start+size-1, size)
	}

	var sb strings.Builder
	for i := range start + size {
		if i < start {
			sb.WriteRune(' ')
		} else if i == start {
			sb.WriteRune('|')
		} else if i == (start + size - 1) {
			sb.WriteRune('|')
			sb.WriteString(fmt.Sprintf(" %d - %d (%d)", start, start+size-1, size))
		} else {
			sb.WriteRune('-')
		}
	}
	return sb.String()
}
