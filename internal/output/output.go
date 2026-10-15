package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
  "strconv"
	"text/tabwriter"
	"time"

	"github.com/CKnuchel/talon/internal/scanner"
)

// PrintTable writes a human-readable table of results to w, followed by
// a one-line summary (ports scanned, open count, elapsed duration).
func PrintTable(w io.Writer, results []scanner.PortResult, elapsed time.Duration) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	// write header
	if _, err := fmt.Fprintf(tw, "PORT\tSTATE\n"); err != nil {
		return err
	}

	open := 0
	for _, r := range results {
		if r.State == scanner.StateOpen {
			open++
		}
		if _, err := fmt.Fprintf(tw, "%d\t%v\n", r.Port, r.State); err != nil {
			return err
		}
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	// summary
	_, err := fmt.Fprintf(w, "\nScanned %d ports in %s - %d open\n", len(results), elapsed, open)
	return err
}

// ExportJSON writes results to w as a JSON array.
func ExportJSON(w io.Writer, results []scanner.PortResult) error {
  return json.NewEncoder(w).Encode(results)
}

// ExportCSV writes results to w as CSV with a header row (port,state).
func ExportCSV(w io.Writer, results []scanner.PortResult) error {
	cw := csv.NewWriter(w)

	// write header
	if err := cw.Write([]string{"port", "state"}); err != nil {
		return err
	}

	for _, r := range results {
		row := []string{strconv.Itoa(r.Port), string(r.State)}
		if err := cw.Write(row); err != nil {
			return err
		}
	}

	cw.Flush()

	return cw.Error()
}