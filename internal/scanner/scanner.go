package scanner

import(
  "strconv"
  "net"
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