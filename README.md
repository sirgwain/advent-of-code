# advent-of-code
This repo contains my solutions to the https://adventofcode.com daily puzzles.

## building

```zsh
make
```

## inputs
inputs are downloaded automatically if an `AOC_SESSION` env var is set. You can also create a `.env` file for the environment variable

```zsh
export AOC_SESSION=53616c7....

./advent-of-code run --year 2015 --day 1

inputs/2015/day1.txt not found, downloading from adventofcode.com
2015-1 solution1: 74 solution2: 1795
Time taken 516µs
```

## running

puzzles can be run for a year/day, or all puzzles can be run if no year is specified

```zsh
❯ ./advent-of-code run -h
run the solution for a year/day or all days if no year/day is specified

Usage:
  advent-of-code-2025 run [flags]

Flags:
  -d, --day int        the puzzle day to run
  -h, --help           help for run
  -i, --input string   the input file to load, defaults to inputs/<year>/day<n>.txt
  -y, --year int       the year to run (default 2015)

Global Flags:
      --debug   enable debug logging
```
