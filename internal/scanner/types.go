package scanner

// PortState represents the state of a scanned port.
type PortState string

const (
	// StateOpen indicates the port accepted a connection.
	StateOpen PortState = "open"
	// StateClosed indicates the port did not accept a connection within the timeout.
	StateClosed PortState = "closed"
)

// PortResult represents the outcome of scanning a single port.
type PortResult struct {
	Port  int       `json:"port"`
	State PortState `json:"state"`
}