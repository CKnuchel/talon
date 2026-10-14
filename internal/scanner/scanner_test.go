package scanner

import (
	"net"
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