# appdash-tools
A set of developer tools for github.com/sourcegraph/appdash

#### Usage
1. From appdash export the trace you'd like to analyze.
2. Create a file `trace.json` within the cloned repo: `.../chris-ramon/appdash-tools/trace.json`
3. Execute `go run main.go` will print sub-traces with/without parent found on the trace root.
