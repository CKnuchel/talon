package output

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/CKnuchel/talon/internal/scanner"
)

func TestExportJSON(t *testing.T) {
	results := []scanner.PortResult{
		{Port: 22, State: scanner.StateOpen},
		{Port: 80, State: scanner.StateOpen},
	}

	var buf bytes.Buffer
	if err := ExportJSON(&buf, results); err != nil {
		t.Fatalf("ExportJSON returned error: %v", err)
	}

	var got []scanner.PortResult
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}

	if !reflect.DeepEqual(got, results) {
		t.Errorf("ExportJSON roundtrip = %v, want %v", got, results)
	}
}

func TestExportCSV(t *testing.T) {
	results := []scanner.PortResult{
		{Port: 22, State: scanner.StateOpen},
		{Port: 443, State: scanner.StateClosed},
	}

	var buf bytes.Buffer
	if err := ExportCSV(&buf, results); err != nil {
		t.Fatalf("ExportCSV returned error: %v", err)
	}

	reader := csv.NewReader(&buf)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("output is not valid CSV: %v", err)
	}

	expected := [][]string{
		{"port", "state"},
		{"22", "open"},
		{"443", "closed"},
	}
	if !reflect.DeepEqual(records, expected) {
		t.Errorf("ExportCSV = %v, want %v", records, expected)
	}
}

func TestPrintTable(t *testing.T) {
	results := []scanner.PortResult{
		{Port: 22, State: scanner.StateOpen},
	}

	var buf bytes.Buffer
	if err := PrintTable(&buf, results, 1200*time.Millisecond); err != nil {
		t.Fatalf("PrintTable returned error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "22") || !strings.Contains(output, "open") {
		t.Errorf("table output missing expected content: %s", output)
	}
}