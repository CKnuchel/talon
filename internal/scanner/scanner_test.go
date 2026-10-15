package scanner

import (
  "fmt"
	"net"
  "reflect"
  "sort"
	"testing"
	"time"
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

	tests := []struct {
		name          string
		host          string
		start         int
		end           int
		includeClosed bool
		expected      []PortResult
	}{
		{
			name:          "open only, includeClosed false",
			host:          "127.0.0.1",
			start:         ipv4Ports[0],
			end:           ipv4Ports[len(ipv4Ports)-1],
			includeClosed: false,
			expected: []PortResult{
				{Port: ipv4Ports[0], State: StateOpen},
				{Port: ipv4Ports[1], State: StateOpen},
				{Port: ipv4Ports[2], State: StateOpen},
			},
		},
		{
			name:          "gap port included as closed",
			host:          "127.0.0.1",
			start:         ipv4Ports[0] - 1,
			end:           ipv4Ports[0],
			includeClosed: true,
			expected: []PortResult{
				{Port: ipv4Ports[0] - 1, State: StateClosed},
				{Port: ipv4Ports[0], State: StateOpen},
			},
		},
	}

	for _, rc := range tests {
		t.Run(rc.name, func(t *testing.T) {
			got := ScanRange(rc.host, rc.start, rc.end, 200*time.Millisecond, rc.includeClosed)
			if !reflect.DeepEqual(got, rc.expected) {
				t.Errorf("ScanRange(...) = %v, want %v", got, rc.expected)
			}
		})
	}
}

func TestScanRangeConcurrent(t *testing.T) {
	ipv4Ports := []int{49501, 49502, 49503, 49504}
	var ipv4Listeners []net.Listener
	for _, p := range ipv4Ports {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", p))
		if err != nil {
			t.Fatalf("failed to start listener on port %d: %v", p, err)
		}
		defer ln.Close()
		ipv4Listeners = append(ipv4Listeners, ln)
	}

	tests := []struct {
		name          string
		host          string
		start         int
		end           int
		workers       int
		includeClosed bool
		expected      []PortResult
	}{
		{
			name:          "open only, multiple workers",
			host:          "127.0.0.1",
			start:         ipv4Ports[0] - 1,
			end:           ipv4Ports[len(ipv4Ports)-1],
			workers:       4,
			includeClosed: false,
			expected: []PortResult{
				{Port: ipv4Ports[0], State: StateOpen},
				{Port: ipv4Ports[1], State: StateOpen},
				{Port: ipv4Ports[2], State: StateOpen},
				{Port: ipv4Ports[3], State: StateOpen},
			},
		},
		{
			name:          "includeClosed true, mixed states",
			host:          "127.0.0.1",
			start:         ipv4Ports[0] - 1,
			end:           ipv4Ports[0],
			workers:       4,
			includeClosed: true,
			expected: []PortResult{
				{Port: ipv4Ports[0] - 1, State: StateClosed},
				{Port: ipv4Ports[0], State: StateOpen},
			},
		},
	}

	for _, rc := range tests {
		t.Run(rc.name, func(t *testing.T) {
			got := ScanRangeConcurrent(rc.host, rc.start, rc.end, 200*time.Millisecond, rc.workers, rc.includeClosed)
			sort.Slice(got, func(i, j int) bool { return got[i].Port < got[j].Port })
			if !reflect.DeepEqual(got, rc.expected) {
				t.Errorf("ScanRangeConcurrent(...) = %v, want %v", got, rc.expected)
			}
		})
	}
}