package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sirgwain/advent-of-code/advent"
	"github.com/spf13/cobra"
)

func runDay(r dayRunner, input string, year, day int) error {
	// load the input data
	defaultInput := fmt.Sprintf("inputs/%d/day%d.txt", year, day)
	if input == "" {
		input = defaultInput
	}

	// if the input file doesn't exist, download it if it's the default day input
	if _, err := os.Stat(input); errors.Is(err, os.ErrNotExist) {
		if input != defaultInput {
			return fmt.Errorf("input %s not found", input)
		}
		fmt.Printf("%s not found, downloading from adventofcode.com\n", input)
		if err := downloadInput(year, day, input); err != nil {
			return err
		}
	}

	content, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file %s %v", input, err)
	}

	start := time.Now()
	defer func() { fmt.Printf("\nTime taken %v\n", time.Since(start)) }()

	// run the day
	s1, s2, err := r.Run(content)
	if err != nil {
		return err
	}

	fmt.Printf("%d-%d ", year, day)
	advent.PrintSolution(s1, s2)

	return nil
}

func newRunCmd() *cobra.Command {
	var year int
	var day int
	var input string
	cmd := &cobra.Command{
		Use:   "run",
		Short: "run a solution",
		Long:  `run the solution for a year/day or all days if no year/day is specified`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var runners []dayRunner
			if day != 0 {
				// get the runner
				r, err := getRunner(year, day)
				if err != nil {
					return err
				}
				runners = append(runners, r)
			} else {
				runners = getAllRunners()
			}

			start := time.Now()
			defer func() {
				if len(runners) > 1 {
					fmt.Printf("\nTotal time taken %v\n", time.Since(start))
				}
			}()

			// run the runners
			for _, d := range runners {
				if err := runDay(d, input, d.year, d.day); err != nil {
					return fmt.Errorf("failed to run %d/%d %v", d.year, d.day, err)
				}
			}

			return nil
		},
	}

	cmd.Flags().IntVarP(&year, "year", "y", 2015, "the year to run")
	cmd.Flags().IntVarP(&day, "day", "d", 0, "the puzzle day to run")
	cmd.Flags().StringVarP(&input, "input", "i", "", "the input file to load, defaults to inputs/<year>/day<n>.txt")

	return cmd
}

func init() {
	rootCmd.AddCommand(newRunCmd())
}
