package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirgwain/advent-of-code/advent"
	"github.com/spf13/cobra"
)

func defaultInputPath(year, day int) string {
	return fmt.Sprintf("inputs/%d/day%d.txt", year, day)
}

// solveDay loads the input (downloading the default input if missing) and
// runs the day, returning both answers and the solve time.
func solveDay(r dayRunner, input string, year, day int) (string, string, time.Duration, error) {
	defaultInput := defaultInputPath(year, day)
	if input == "" {
		input = defaultInput
	}

	// if the input file doesn't exist, download it if it's the default day input
	if _, err := os.Stat(input); errors.Is(err, os.ErrNotExist) {
		if input != defaultInput {
			return "", "", 0, fmt.Errorf("input %s not found", input)
		}
		fmt.Printf("%s not found, downloading from adventofcode.com\n", input)
		if err := downloadInput(year, day, input); err != nil {
			return "", "", 0, err
		}
	}

	content, err := os.ReadFile(input)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to read input file %s %v", input, err)
	}

	start := time.Now()
	s1, s2, err := r.Run(content)
	elapsed := time.Since(start)
	return s1, s2, elapsed, err
}

// runDay runs a day and prints its answers. When running against the default
// input, answers are compared to inputs/solutions.json; it returns true if
// they differ from the recorded solution.
func runDay(r dayRunner, input string, sols solutions, year, day int) (bool, error) {
	s1, s2, elapsed, err := solveDay(r, input, year, day)
	if err != nil {
		return false, err
	}

	fmt.Printf("%d-%d ", year, day)
	advent.PrintSolution(s1, s2)

	changed := false
	if input == "" || input == defaultInputPath(year, day) {
		var status string
		status, changed = sols.check(year, day, s1, s2)
		fmt.Printf(" %s", status)
	}
	fmt.Printf("\nTime taken %v\n", elapsed)

	return changed, nil
}

// selectRunners returns the runner for year/day, or every runner if day is 0.
func selectRunners(year, day int) ([]dayRunner, error) {
	if day == 0 {
		return getAllRunners(), nil
	}
	r, err := getRunner(year, day)
	if err != nil {
		return nil, err
	}
	return []dayRunner{r}, nil
}

func newRunCmd() *cobra.Command {
	var year int
	var day int
	var input string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a solution",
		Long: `run the solution for a year/day or all days if no year/day is specified.
Answers are checked against inputs/solutions.json (see the record command).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			runners, err := selectRunners(year, day)
			if err != nil {
				return err
			}

			sols, err := loadSolutions()
			if err != nil {
				return err
			}

			start := time.Now()
			defer func() {
				if len(runners) > 1 {
					fmt.Printf("\nTotal time taken %v\n", time.Since(start))
				}
			}()

			// run the runners
			var changed []string
			for _, d := range runners {
				c, err := runDay(d, input, sols, d.year, d.day)
				if err != nil {
					return fmt.Errorf("failed to run %d/%d %v", d.year, d.day, err)
				}
				if c {
					changed = append(changed, solutionKey(d.year, d.day))
				}
			}

			if len(changed) > 0 {
				cmd.SilenceUsage = true
				return fmt.Errorf("%d solution(s) differ from %s: %v", len(changed), solutionsFile, changed)
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&year, "year", "y", 2015, "the year to run")
	cmd.Flags().IntVarP(&day, "day", "d", 0, "the puzzle day to run")
	cmd.Flags().BoolVar(&advent.Viz.Enabled, "viz", false, "enable visualization")
	cmd.Flags().StringVarP(&input, "input", "i", "", "the input file to load, defaults to inputs/<year>/day<n>.txt")

	return cmd
}

func newRecordCmd() *cobra.Command {
	var year int
	var day int
	cmd := &cobra.Command{
		Use:   "record",
		Short: "record solutions to inputs/solutions.json",
		Long: `run the solution for a year/day (or all days if no day is specified) against
the default input and record the answers in inputs/solutions.json. Subsequent
runs will flag any answer that differs from the recorded one.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			runners, err := selectRunners(year, day)
			if err != nil {
				return err
			}

			sols, err := loadSolutions()
			if err != nil {
				return err
			}

			for _, d := range runners {
				s1, s2, _, err := solveDay(d, "", d.year, d.day)
				if err != nil {
					return fmt.Errorf("failed to run %d/%d %v", d.year, d.day, err)
				}

				key := solutionKey(d.year, d.day)
				fmt.Printf("%d-%d ", d.year, d.day)
				advent.PrintSolution(s1, s2)
				if old, ok := sols[key]; ok && (old.Part1 != s1 || old.Part2 != s2) {
					fmt.Printf(" %s", mismatchStyle.Render(fmt.Sprintf("(was %s / %s)", old.Part1, old.Part2)))
				}
				fmt.Println()

				sols[key] = solution{Part1: s1, Part2: s2}
			}

			if err := sols.save(); err != nil {
				return err
			}
			fmt.Printf("recorded %d solution(s) to %s\n", len(runners), solutionsFile)
			return nil
		},
	}

	cmd.Flags().IntVarP(&year, "year", "y", 2015, "the year to record")
	cmd.Flags().IntVarP(&day, "day", "d", 0, "the puzzle day to record (all days if omitted)")

	return cmd
}

func init() {
	rootCmd.AddCommand(newRunCmd())
	rootCmd.AddCommand(newRecordCmd())
}
