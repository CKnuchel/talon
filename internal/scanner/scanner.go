package scanner

import (
	"net"
	"strconv"
	"sync"
	"time"
)

// ScanPort checks whether a single TCP port is open on the given host.
// It returns true if a connection succeeds within the given timeout.
func ScanPort(host string, port int, timeout time.Duration) bool {
	target := net.JoinHostPort(host, strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()

	return true
}

// ScanRange scans all ports in the range [start, end] on the given host.
// If includeClosed is true, closed ports are included in the result with
// StateClosed; otherwise only open ports are returned.
func ScanRange(host string, start, end int, timeout time.Duration, includeClosed bool) []PortResult {
	ports := []PortResult{}

	for i := start; i <= end; i++ {
		open := ScanPort(host, i, timeout)

		if !open && !includeClosed {
			continue
		}

		state := StateClosed
		if open {
			state = StateOpen
		}

		ports = append(ports, PortResult{Port: i, State: state})
	}

	return ports
}

// ScanRangeConcurrent scans all ports in the range [start, end] on the given host
// using multiple workers in parallel. If includeClosed is true, closed ports are
// included in the result with StateClosed; otherwise only open ports are returned.
func ScanRangeConcurrent(host string, start, end int, timeout time.Duration, workers int, includeClosed bool) []PortResult {
  if start > end {
		return []PortResult{}
  }
  
	maxJobs := end - start + 1
	jobs := make(chan int, maxJobs)
	result := make(chan PortResult, maxJobs)
	var wg sync.WaitGroup

	// fill jobs channel
	for i := start; i <= end; i++ {
		jobs <- i
	}
	close(jobs)

	// start workers
	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go worker(host, jobs, result, timeout, &wg, includeClosed)
	}

	// wait for all workers and close result chan
	go func() {
		wg.Wait()
		close(result)
	}()

	// collect results
	ports := []PortResult{}
	for res := range result {
		ports = append(ports, res)
	}

	return ports
}

func worker(host string, jobs chan int, result chan PortResult, timeout time.Duration, wg *sync.WaitGroup, includeClosed bool) {
	defer wg.Done()

	for port := range jobs {
		open := ScanPort(host, port, timeout)

		if !open && !includeClosed {
			continue
		}

		state := StateClosed
		if open {
			state = StateOpen
		}

		result <- PortResult{Port: port, State: state}
	}
}