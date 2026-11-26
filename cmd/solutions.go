package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/sirgwain/advent-of-code/advent/color"
)

// solutionsFile holds known-good answers so `run` can flag regressions while
// optimizing. Populate it with `advent-of-code record`.
const solutionsFile = "inputs/solutions.json"

type solution struct {
	Part1 string `json:"part1"`
	Part2 string `json:"part2"`
}

// solutions is keyed by "<year>/dayNN" so encoding/json writes it sorted.
type solutions map[string]solution

func solutionKey(year, day int) string {
	return fmt.Sprintf("%d/day%02d", year, day)
}

func loadSolutions() (solutions, error) {
	sols := solutions{}
	data, err := os.ReadFile(solutionsFile)
	if errors.Is(err, os.ErrNotExist) {
		return sols, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &sols); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", solutionsFile, err)
	}
	return sols, nil
}

func (s solutions) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(solutionsFile), 0o755); err != nil {
		return err
	}
	return os.WriteFile(solutionsFile, append(data, '\n'), 0o644)
}

var (
	matchStyle    = lipgloss.NewStyle().Foreground(color.BrightGreen82)
	mismatchStyle = lipgloss.NewStyle().Foreground(color.LightRed9).Bold(true)
	unknownStyle  = lipgloss.NewStyle().Faint(true)
)

// checkSolution returns a short status string comparing s1/s2 to the
// recorded answers for year/day, and whether they differ.
func (s solutions) check(year, day int, s1, s2 string) (string, bool) {
	want, ok := s[solutionKey(year, day)]
	if !ok {
		return unknownStyle.Render("(no recorded solution)"), false
	}
	if want.Part1 == s1 && want.Part2 == s2 {
		return matchStyle.Render("✓"), false
	}
	msg := "✗ CHANGED"
	if want.Part1 != s1 {
		msg += fmt.Sprintf(" part1 expected %s", want.Part1)
	}
	if want.Part2 != s2 {
		msg += fmt.Sprintf(" part2 expected %s", want.Part2)
	}
	return mismatchStyle.Render(msg), true
}
