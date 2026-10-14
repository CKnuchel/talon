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