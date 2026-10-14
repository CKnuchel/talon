package output

import (
  "io"
  "time"

  "github.com/CKnuchel/talon/internal/scanner"
)

// PrintTable writes a human-readable table of results to w, followed by
// a one-line summary (ports scanned, open count, elapsed duration).
func PrintTable(w io.Writer, results []scanner.PortResult, elapsed time.Duration) error

// ExportJSON writes results to w as a JSON array.
func ExportJSON(w io.Writer, results []scanner.PortResult) error

// ExportCSV writes results to w as CSV with a header row (port,state).
func ExportCSV(w io.Writer, results []scanner.PortResult) error