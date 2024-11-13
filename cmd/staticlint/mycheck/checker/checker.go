// Package checker defines checker that defines list of analysers including custom, public, static analysers.
//
// For basic usage build binary file in ./../mycheck directory.
//
//	go build -o mycheck main.go
//
// Locate binary file into the project's root.
// To check all packages beneath the current directory:
// mycheck ./...
package checker

import (
	"github.com/Archetarcher/metrics.git/cmd/staticlint/mycheck/analyzers"
	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"
	"slices"
)

// Check is a function that defines and registers all checkers into multichecker
func Check() {
	mychecks := analyzers.StaticAnalyzers

	for _, v := range staticcheck.Analyzers {
		if slices.Contains(analyzers.SAStaticChecks, v.Analyzer.Name) {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
