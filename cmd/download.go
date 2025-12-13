package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func getAOCSession() (string, error) {
	s := os.Getenv("AOC_SESSION")
	if s != "" {
		return s, nil
	}
	if err := loadDotEnv(".env"); err != nil {
		return "", err
	}
	s = os.Getenv("AOC_SESSION")
	if s == "" {
		return "", fmt.Errorf("AOC_SESSION not set (env or .env)")
	}
	return s, nil
}

func downloadInput(year, day int, outPath string) error {
	if day < 1 || day > 25 {
		return fmt.Errorf("day must be 1..25")
	}

	session, err := getAOCSession()
	if err != nil {
		return err
	}

	// Make parent dirs
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	// If file exists, don’t redownload
	if _, err := os.Stat(outPath); err == nil {
		return nil
	}

	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Cookie", "session="+session)
	req.Header.Set("User-Agent", "github.com/sirgwain/advent-of-code")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("download failed: %s: %s", resp.Status, string(b))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return os.WriteFile(outPath, body, 0o644)
}

func newDownloadCmd() *cobra.Command {
	var year int
	var day int
	var outDir string
	var force bool

	cmd := &cobra.Command{
		Use:   "download",
		Short: "download Advent of Code input",
		Long:  "Download your personal Advent of Code puzzle input",
		RunE: func(cmd *cobra.Command, args []string) error {
			if day < 1 || day > 25 {
				return fmt.Errorf("day must be between 1 and 25")
			}
			outPath := filepath.Join(
				outDir,
				strconv.Itoa(year),
				fmt.Sprintf("day%d.txt", day),
			)

			if !force {
				if _, err := os.Stat(outPath); err == nil {
					fmt.Printf("input already exists: %s\n", outPath)
					return nil
				}
			}

			if err := downloadInput(year, day, outPath); err != nil {
				return err
			}

			fmt.Printf("downloaded input → %s\n", outPath)
			return nil
		},
	}

	cmd.Flags().IntVarP(&year, "year", "y", time.Now().Year(), "AoC year")
	cmd.Flags().IntVarP(&day, "day", "d", 0, "AoC day (1-25)")
	cmd.Flags().StringVarP(&outDir, "out", "o", "inputs", "output directory")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "overwrite existing file")

	cmd.MarkFlagRequired("day")

	return cmd
}

func init() {
	rootCmd.AddCommand(newDownloadCmd())
}
