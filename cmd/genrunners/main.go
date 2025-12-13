package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"go/token"
	"go/types"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type entry struct {
	Year       int
	Day        int
	Alias      string // import alias e.g. y2015d1
	ImportPath string
}

var (
	outFile = flag.String("out", "cmd/runner_gen.go", "output file")
	module  = flag.String("module", "", "module import path prefix, e.g. github.com/sirgwain/advent-of-code")
	rootDir = flag.String("root", ".", "filesystem root to scan")
	pkgBase = flag.String("pkgbase", "advent", "base directory containing year/day packages (relative to root)")
)

func main() {
	flag.Parse()

	if *module == "" {
		fatal(errors.New(`-module is required, e.g. -module github.com/sirgwain/advent-of-code`))
	}

	entries, err := discoverEntries(*rootDir, *pkgBase, *module)
	if err != nil {
		fatal(err)
	}
	if len(entries) == 0 {
		fatal(fmt.Errorf("no valid day packages found under %s/%s", *rootDir, *pkgBase))
	}

	src, err := render(entries)
	if err != nil {
		fatal(err)
	}

	// gofmt
	formatted, err := format.Source(src)
	if err != nil {
		// write raw for debugging if formatting fails
		_ = os.WriteFile(*outFile, src, 0644)
		fatal(fmt.Errorf("gofmt failed: %w (wrote unformatted output to %s)", err, *outFile))
	}

	if err := os.MkdirAll(filepath.Dir(*outFile), 0755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(*outFile, formatted, 0644); err != nil {
		fatal(err)
	}
	fmt.Printf("generated %s (%d runners)\n", *outFile, len(entries))
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}

// discoverEntries finds advent/<year>/day<day> directories, type-checks them,
// and returns entries where *Day implements: Run([]byte) (int,int,error)
func discoverEntries(root, pkgbase, module string) ([]entry, error) {
	yearDayRE := regexp.MustCompile(`^(\d{4})[/\\]day(\d+)$`)

	var candidates []struct {
		year int
		day  int
		dir  string // filesystem dir
		imp  string // import path
	}

	basePath := filepath.Join(root, pkgbase)

	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(basePath, path)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)

		m := yearDayRE.FindStringSubmatch(rel)
		if m == nil {
			return nil
		}

		year := atoi(m[1])
		day := atoi(m[2])

		imp := strings.TrimRight(module, "/") + "/" + filepath.ToSlash(filepath.Join(pkgbase, fmt.Sprintf("%04d", year), fmt.Sprintf("day%d", day)))

		candidates = append(candidates, struct {
			year int
			day  int
			dir  string
			imp  string
		}{year: year, day: day, dir: path, imp: imp})

		return nil
	})
	if err != nil {
		return nil, err
	}

	// We type-check by import path (more reliable for go/packages), but keep dirs for debug if you want.
	var entries []entry
	for _, c := range candidates {
		ok, err := packageHasRunnableDay(c.imp)
		if err != nil {
			// treat load/type issues as "not a runner"; comment out to fail hard
			continue
		}
		if !ok {
			continue
		}

		alias := fmt.Sprintf("y%dd%d", c.year, c.day)
		entries = append(entries, entry{
			Year:       c.year,
			Day:        c.day,
			Alias:      alias,
			ImportPath: c.imp,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Year != entries[j].Year {
			return entries[i].Year < entries[j].Year
		}
		return entries[i].Day < entries[j].Day
	})

	return entries, nil
}

func packageHasRunnableDay(importPath string) (bool, error) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax | packages.NeedName | packages.NeedFiles,
	}
	pkgs, err := packages.Load(cfg, importPath)
	if err != nil {
		return false, err
	}
	if packages.PrintErrors(pkgs) > 0 {
		return false, fmt.Errorf("package load errors for %s", importPath)
	}
	if len(pkgs) == 0 || pkgs[0].Types == nil {
		return false, fmt.Errorf("no types for %s", importPath)
	}
	p := pkgs[0].Types

	// Build the runner interface type: interface{ Run([]byte) (int,int,error) }
	// We'll create it as a *types.Interface so we can use types.Implements.
	byteSlice := types.NewSlice(types.Typ[types.Byte])

	runSig := types.NewSignatureType(
		nil,
		nil, nil,
		types.NewTuple(types.NewVar(token.NoPos, nil, "", byteSlice)),
		types.NewTuple(
			types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]),
			types.NewVar(token.NoPos, nil, "", types.Typ[types.Int]),
			types.NewVar(token.NoPos, nil, "", types.Universe.Lookup("error").Type()),
		),
		false,
	)

	runFn := types.NewFunc(token.NoPos, nil, "Run", runSig)
	runnerIface := types.NewInterfaceType([]*types.Func{runFn}, nil)
	runnerIface.Complete()

	// Look for exported type "Day"
	obj := p.Scope().Lookup("Day")
	if obj == nil {
		return false, nil
	}
	tn, ok := obj.(*types.TypeName)
	if !ok {
		return false, nil
	}

	// We want *Day to implement runner.
	ptr := types.NewPointer(tn.Type())
	return types.Implements(ptr, runnerIface), nil
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			break
		}
		n = n*10 + int(r-'0')
	}
	return n
}

func render(entries []entry) ([]byte, error) {
	// Group by year for nested switches.
	type yearGroup struct {
		Year int
		Days []entry
	}
	var groups []yearGroup
	for i := 0; i < len(entries); {
		y := entries[i].Year
		j := i
		for j < len(entries) && entries[j].Year == y {
			j++
		}
		groups = append(groups, yearGroup{Year: y, Days: entries[i:j]})
		i = j
	}

	// Deterministic import ordering: by ImportPath.
	imports := make([]entry, len(entries))
	copy(imports, entries)
	sort.Slice(imports, func(i, j int) bool { return imports[i].ImportPath < imports[j].ImportPath })

	data := struct {
		Imports []entry
		Groups  []yearGroup
		Total   int
	}{
		Imports: imports,
		Groups:  groups,
		Total:   len(entries),
	}

	tmpl := template.Must(template.New("gen").Parse(genTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

const genTemplate = `// Code generated by cmd/genrunners; DO NOT EDIT.

package cmd

import (
	"fmt"
{{- range .Imports }}
	{{ .Alias }} "{{ .ImportPath }}"
{{- end }}
)

type runner interface {
	Run([]byte) (int, int, error)
}

type dayRunner struct {
	r    runner
	year int
	day  int
}

func getRunner(year, day int) (dayRunner, error) {
	var r runner
	switch year {
{{- range .Groups }}
	case {{ .Year }}:
		switch day {
		{{- range .Days }}
		case {{ .Day }}:
			r = &{{ .Alias }}.Day{}
		{{- end }}
		}
{{- end }}
	}

	if r == nil {
		return dayRunner{}, fmt.Errorf("day %d/%d not found", year, day)
	}
	return dayRunner{r: r, year: year, day: day}, nil
}

func getAllRunners() ([]dayRunner) {
	runners := make([]dayRunner, 0, {{.Total}})
{{- range .Groups }}
	{{- range .Days }}
	runners = append(runners, dayRunner{r: &{{ .Alias }}.Day{}, year: {{ .Year }}, day: {{ .Day }}})
	{{- end }}
{{- end }}

	return runners
}
`
