package scanner

import (
  "fmt"
	"net"
	"testing"
	"time"
  "reflect"
)

func TestScanPort(t *testing.T) {
	// Set up an open port on IPv4: a real listener accepting connections.
	openLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start open-port listener: %v", err)
	}
	defer openLn.Close()
	openPort := openLn.Addr().(*net.TCPAddr).Port

	// Set up an open port on IPv6 as well.
	openLn6, err := net.Listen("tcp", "[::1]:0")
	if err != nil {
		t.Fatalf("failed to start open-port listener (IPv6): %v", err)
	}
	defer openLn6.Close()
	openPort6 := openLn6.Addr().(*net.TCPAddr).Port

	// Set up a closed port: open a listener, grab its port, then close it.
	closedLn, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start closed-port listener: %v", err)
	}
	closedPort := closedLn.Addr().(*net.TCPAddr).Port
	closedLn.Close()

	tests := []struct {
		name     string
		host     string
		port     int
		expected bool
	}{
		{name: "open port (IPv4)", host: "127.0.0.1", port: openPort, expected: true},
		{name: "open port (IPv6)", host: "::1", port: openPort6, expected: true},
		{name: "closed port", host: "127.0.0.1", port: closedPort, expected: false},
	}

	for _, rc := range tests {
		t.Run(rc.name, func(t *testing.T) {
			got := ScanPort(rc.host, rc.port, 200*time.Millisecond)
			if got != rc.expected {
				t.Errorf("ScanPort(%s, %d) = %v, want %v", rc.host, rc.port, got, rc.expected)
			}
		})
	}
}

func TestScanRange(t *testing.T) {
	// Fixed, closely-spaced ports in a high, unlikely-to-be-used range.
	ipv4Ports := []int{49201, 49202, 49203}
	var ipv4Listeners []net.Listener
	for _, p := range ipv4Ports {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err != nil {
			t.Fatalf("failed to start listener on port %d: %v", p, err)
		}
		defer ln.Close()
		ipv4Listeners = append(ipv4Listeners, ln)
	}

	ipv6Ports := []int{49301, 49302}
	var ipv6Listeners []net.Listener
	for _, p := range ipv6Ports {
		ln, err := net.Listen("tcp", fmt.Sprintf("[::1]:%d", p))
		if err != nil {
			t.Fatalf("failed to start listener on port %d (IPv6): %v", p, err)
		}
		defer ln.Close()
		ipv6Listeners = append(ipv6Listeners, ln)
	}

	// Single port used to test boundaries (start == end).
	boundaryPort := 49401
	boundaryLn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", boundaryPort))
	if err != nil {
		t.Fatalf("failed to start listener on port %d: %v", boundaryPort, err)
	}
	defer boundaryLn.Close()

	tests := []struct {
		name     string
		host     string
		start    int
		end      int
		expected []int
	}{
		{
			name:     "IPv4 range with open ports",
			host:     "127.0.0.1",
			start:    ipv4Ports[0] - 1,
			end:      ipv4Ports[len(ipv4Ports)-1] + 1,
			expected: ipv4Ports,
		},
		{
			name:     "IPv6 range with open ports",
			host:     "::1",
			start:    ipv6Ports[0] - 1,
			end:      ipv6Ports[len(ipv6Ports)-1] + 1,
			expected: ipv6Ports,
		},
		{
			name:     "start equals end, port open",
			host:     "127.0.0.1",
			start:    boundaryPort,
			end:      boundaryPort,
			expected: []int{boundaryPort},
		},
		{
			name:     "start equals end, port closed",
			host:     "127.0.0.1",
			start:    boundaryPort + 1,
			end:      boundaryPort + 1,
			expected: []int{},
		},
		{
			name:     "end port is inclusive",
			host:     "127.0.0.1",
			start:    boundaryPort - 1,
			end:      boundaryPort,
			expected: []int{boundaryPort},
		},
		{
			name:     "start greater than end returns empty",
			host:     "127.0.0.1",
			start:    boundaryPort,
			end:      boundaryPort - 1,
			expected: []int{},
		},
	}

	for _, rc := range tests {
		t.Run(rc.name, func(t *testing.T) {
			got := ScanRange(rc.host, rc.start, rc.end, 200*time.Millisecond)
			if !reflect.DeepEqual(got, rc.expected) {
				t.Errorf("ScanRange(%s, %d, %d) = %v, want %v", rc.host, rc.start, rc.end, got, rc.expected)
			}
		})
	}
}