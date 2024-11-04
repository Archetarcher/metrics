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
	"github.com/Archetarcher/metrics.git/cmd/staticlint/osexitchecker"
	"github.com/go-critic/go-critic/checkers/analyzer"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/analysis/facts/directives"
	"honnef.co/go/tools/analysis/facts/nilness"
	"honnef.co/go/tools/staticcheck"
	"slices"
)

// StaticAnalyzers is a slice of analyzers from golang.org/x/tools/go/analysis/passes
var StaticAnalyzers = []*analysis.Analyzer{
	appends.Analyzer,         // detects if there is only one variable in append.
	asmdecl.Analyzer,         // reports mismatches between assembly files and Go declarations.
	assign.Analyzer,          // detects useless assignments.
	atomic.Analyzer,          // checks for common mistakes using the sync/atomic package.
	atomicalign.Analyzer,     // checks for non-64-bit-aligned arguments to sync/atomic functions.
	bools.Analyzer,           // detects common mistakes involving boolean operators.
	buildssa.Analyzer,        // constructs the SSA representation of an error-free package and returns the set of all functions within it.
	buildtag.Analyzer,        // checks build tags.
	cgocall.Analyzer,         // detects some violations of the cgo pointer passing rules.
	composite.Analyzer,       // checks for unkeyed composite literals.
	copylock.Analyzer,        // checks for locks erroneously passed by value.
	ctrlflow.Analyzer,        // provides a syntactic control-flow graph (CFG) for the body of a function.
	deepequalerrors.Analyzer, // checks for the use of reflect.DeepEqual with error values.
	defers.Analyzer,          // checks for common mistakes in defer statements.
	directives.Analyzer,      // checks known Go toolchain directives.
	errorsas.Analyzer,        // checks that the second argument to errors.As is a pointer to a type implementing error.
	fieldalignment.Analyzer,  // detects structs that would use less memory if their fields were sorted.
	findcall.Analyzer,        // serves as a trivial example and test of the Analysis API.
	framepointer.Analyzer,    // reports assembly code that clobbers the frame pointer before saving it.
	httpmux.Analyzer,
	httpresponse.Analyzer,        // checks for mistakes using HTTP responses.
	ifaceassert.Analyzer,         // flags impossible interface-interface type assertions.
	inspect.Analyzer,             // provides an AST inspector (golang.org/x/tools/go/ast/inspector.Inspector) for the syntax trees of a package.
	loopclosure.Analyzer,         // checks for references to enclosing loop variables from within nested functions.
	lostcancel.Analyzer,          // checks for failure to call a context cancellation function.
	nilfunc.Analyzer,             // checks for useless comparisons against nil.
	nilness.Analysis,             // inspects the control-flow graph of an SSA function and reports errors such as nil pointer dereferences and degenerate nil pointer comparisons.
	pkgfact.Analyzer,             // is a demonstration and test of the package fact mechanism.
	printf.Analyzer,              // checks consistency of Printf format strings and arguments.
	reflectvaluecompare.Analyzer, // checks for accidentally using == or reflect.DeepEqual to compare reflect.Value values.
	shadow.Analyzer,              // checks for shadowed variables.
	shift.Analyzer,               // checks for shifts that exceed the width of an integer.
	sigchanyzer.Analyzer,         // detects misuse of unbuffered signal as argument to signal.Notify.
	slog.Analyzer,                // for mismatched key-value pairs in log/slog calls.
	sortslice.Analyzer,           // checks for calls to sort.Slice that do not use a slice type as first argument.
	stdmethods.Analyzer,          // checks for misspellings in the signatures of methods similar to well-known interfaces.
	stdversion.Analyzer,          // reports uses of standard library symbols that are "too new" for the Go version in force in the referring file.
	stringintconv.Analyzer,       // flags type conversions from integers to strings.
	structtag.Analyzer,           // checks struct field tags are well formed.
	testinggoroutine.Analyzer,    // for detecting calls to Fatal from a test goroutine.
	tests.Analyzer,               // checks for common mistaken usages of tests and examples.
	timeformat.Analyzer,          // checks for the use of time.Format or time.Parse calls with a bad format.
	unmarshal.Analyzer,           // checks for passing non-pointer or non-interface types to unmarshal and decode functions.
	unreachable.Analyzer,         // checks for unreachable code.
	unsafeptr.Analyzer,           // checks for invalid conversions of uintptr to unsafe.Pointer.
	unusedresult.Analyzer,        // checks for unused results of calls to certain pure functions.
	unusedwrite.Analyzer,         // checks for unused writes to the elements of a struct or array object.
	usesgenerics.Analyzer,        // defines an Analyzer that checks for usage of generic features added in Go 1.18.
	errcheck.Analyzer,            // errcheck is a program for checking for unchecked errors in Go code.
	analyzer.Analyzer,            // go-critic analyzer is a highly extensible Go source code linter providing checks currently missing from other linters.
	osexitchecker.Analyzer,       // is a custom Analyzer that checks for os.Exit() calls in main package
}

// SAStaticChecks is a slice of SA checks from staticcheck.io package
var SAStaticChecks = []string{
	"SA1000",
	"SA1001",
	"SA1002",
	"SA1003",
	"SA1004",
	"SA1005",
	"SA1006",
	"SA1007",
	"SA1008",
	"SA1010",
	"SA1011",
	"SA1012",
	"SA1013",
	"SA1014",
	"SA1015",
	"SA1016",
	"SA1017",
	"SA1018",
	"SA1019",
	"SA1020",
	"SA1021",
	"SA1023",
	"SA1024",
	"SA1025",
	"SA1026",
	"SA1027",
	"SA1028",
	"SA1029",
	"SA1030",
	"SA1031",
	"SA1032",
	"SA2000",
	"SA2001",
	"SA2002",
	"SA2003",
	"SA3000",
	"SA3001",
	"SA4000",
	"SA4001",
	"SA4003",
	"SA4004",
	"SA4005",
	"SA4006",
	"SA4008",
	"SA4009",
	"SA4010",
	"SA4011",
	"SA4012",
	"SA4013",
	"SA4014",
	"SA4015",
	"SA4016",
	"SA4017",
	"SA4018",
	"SA4019",
	"SA4020",
	"SA4021",
	"SA4022",
	"SA4023",
	"SA4024",
	"SA4025",
	"SA4026",
	"SA4027",
	"SA4028",
	"SA4029",
	"SA4030",
	"SA4031",
	"SA4032",
	"SA5000",
	"SA5001",
	"SA5002",
	"SA5003",
	"SA5004",
	"SA5005",
	"SA5007",
	"SA5008",
	"SA5009",
	"SA5010",
	"SA5011",
	"SA5012",
	"SA6000",
	"SA6001",
	"SA6002",
	"SA6003",
	"SA6005",
	"SA6006",
	"SA9001",
	"SA9002",
	"SA9003",
	"SA9004",
	"SA9005",
	"SA9006",
	"SA9007",
	"SA9008",
	"SA9009",
	"ST1000",
	"ST1005",
}

// Check is a function that defines and registers all checkers into multichecker
func Check() {
	mychecks := StaticAnalyzers

	for _, v := range staticcheck.Analyzers {
		if slices.Contains(SAStaticChecks, v.Analyzer.Name) {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(
		mychecks...,
	)
}
