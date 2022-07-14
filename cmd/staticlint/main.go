// Package staticlint This package is intended for static code analyzation.
package main

import (
	analyzer2 "github.com/antonioo83/shot-url-service/analyzer"
	"github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/gostaticanalysis/nilerr"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)
import _ "net/http/pprof"

// main for check your code using command: checker.exe path_to_filename (checker.exe test.go)
func main() {
	// Define a map of the included rules.
	checks := map[string]bool{
		"SA5000": true, // Assignment to nil map.
		"SA4006": true, // A value assigned to a variable is never read before being overwritten. Forgotten error check or dead code?
		"SA6000": true, // Using regexp.Match or related in a loop, should use regexp.Compile.
		"SA9004": true, // Only the first constant has an explicit type.
		"SA1004": true, // Suspiciously small untyped constant in time.Sleep.
	}
	var mychecks []*analysis.Analyzer
	for _, a := range staticcheck.Analyzers {
		for _, v := range a.Requires {
			if checks[v.Name] {
				mychecks = append(mychecks, v)
			}
		}
	}

	mychecks = append(mychecks, analyzer2.OsExitCheckAnalyzer) // Check the use of the expression "exit" in the "main" function of the "main" package.
	mychecks = append(mychecks, printf.Analyzer)               // Defines an Analyzer that checks consistency of Printf format strings and arguments.
	mychecks = append(mychecks, shadow.Analyzer)               // This analyzer check for shadowed variables.
	mychecks = append(mychecks, structtag.Analyzer)            // Defines an Analyzer that checks struct field tags are well formed.
	mychecks = append(mychecks, analyzer.Analyzer)             // Exports go-critic checkers as analysis-compatible object.
	mychecks = append(mychecks, nilerr.Analyzer)               // Checks returning nil when err is not nil.

	multichecker.Main(mychecks...)
}
