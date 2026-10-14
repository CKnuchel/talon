package main

import (
  "flag"
  "fmt"
  "os"
  "time"

  "github.com/CKnuchel/talon/internal/scanner"
)

func main() {
  // Define allowed flags
  timeout := flag.Duration("t", 200*time.Millisecond, "connection timeout")
  startPort := flag.Int("s", 1, "start port")
  endPort := flag.Int("e", 1024, "end port")

  flag.Parse()

  args := flag.Args()
  if len(args) < 1 {
	  fmt.Println("Usage: talon [flags] <host>")
	  os.Exit(1)
  }
  
  host := args[0]

  // get ports
  openPorts := scanner.ScanRange(host, *startPort, *endPort, *timeout)

  // output result
  for _, port := range openPorts {
    fmt.Printf("Port %d open\n", port)
  }
  
}