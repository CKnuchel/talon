package scanner

import(
  "strconv"
  "net"
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

// ScanRange scans all ports in the range [start, end] on the given host
// and returns the list of ports found to be open.
func ScanRange(host string, start, end int, timeout time.Duration) []int{
  openPorts := []int{}

  for i := start; i <= end; i++ {
    if ScanPort(host, i, timeout) {
      openPorts = append(openPorts, i)
    }
  }

  return openPorts
}

// ScanRangeConcurrent scans all ports in the range [start, end] on the given host
// using multiple workers in parallel and returns the list of ports found to be open.
func ScanRangeConcurrent(host string, start, end int, timeout time.Duration, workers int) []int{
  maxJobs := end - start + 1
  jobs := make(chan int, maxJobs)
  result := make(chan int, maxJobs)
  var wg sync.WaitGroup

  // fill jobs channel
  for i := start; i <= end; i++ {
    jobs <- i
  }
  close(jobs)

  // start workers
  for i := 1; i <= workers; i++ {
    wg.Add(1)
    go worker(host, jobs, result, timeout, &wg)
  }

  // wait for all workers and close result chan
  go func() {
    wg.Wait()
    close(result)
  }()

  // collect results
  openPorts := []int{}
  for op := range result {
    openPorts = append(openPorts, op)
  }
  
  return openPorts
}

func worker(host string, jobs, result chan int, timeout time.Duration, wg *sync.WaitGroup) {
  defer wg.Done()

  for port := range jobs {
    if ScanPort(host, port, timeout) {
      result <- port
    }
  }
}