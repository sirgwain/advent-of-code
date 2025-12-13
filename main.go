//go:generate go run ./cmd/genrunners -module github.com/sirgwain/advent-of-code
package main

import (
	"github.com/sirgwain/advent-of-code/cmd"
)

func main() {
	cmd.Execute()
}
