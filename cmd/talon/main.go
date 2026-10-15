package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/CKnuchel/talon/internal/output"
	"github.com/CKnuchel/talon/internal/scanner"
)

func main() {
	// Define allowed flags
	timeout := flag.Duration("t", 200*time.Millisecond, "connection timeout")
	startPort := flag.Int("s", 1, "start port")
	endPort := flag.Int("e", 1024, "end port")
	workers := flag.Int("w", 100, "number of parallel scans")
	includeClosed := flag.Bool("c", false, "include closed ports in result")
	exportFile := flag.String("o", "", "output file (.json or .csv)")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: talon [flags] <host>")
		os.Exit(1)
	}

	host := args[0]

	// get ports
	start := time.Now()
	result := scanner.ScanRangeConcurrent(host, *startPort, *endPort, *timeout, *workers, *includeClosed)
	elapsed := time.Since(start)

	// ensure sorted output because of goroutines
	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})

	scannedCount := *endPort - *startPort + 1

	// output result
	if err := output.PrintTable(os.Stdout, result, scannedCount, elapsed); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	filename := *exportFile
	if filename == "" {
		return
	}

	// validate extension BEFORE touching the file
	ext := filepath.Ext(filename)
	if ext != ".json" && ext != ".csv" {
		fmt.Println("Error: unsupported export file extension, use .json or .csv")
		os.Exit(1)
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer file.Close()

	switch ext {
	case ".json":
		err = output.ExportJSON(file, result)
	case ".csv":
		err = output.ExportCSV(file, result)
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
