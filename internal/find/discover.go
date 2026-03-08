package find

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"
	"time"
)

const ICSPPort = 1319

// Discover listens on UDP port 1319 for ICSP identify-reply packets and
// returns deduplicated devices after the given timeout, sorted by IP.
func Discover(timeout time.Duration) ([]Device, error) {
	conn, err := listenUDP(ICSPPort)
	if err != nil {
		return nil, err
	}

	defer func() { _ = conn.Close() }()

	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, fmt.Errorf("failed to set read deadline: %w", err)
	}

	return collect(conn, nil), nil
}

// DiscoverWithContext listens on UDP port 1319 until ctx is cancelled (e.g.
// via Ctrl+C), then returns deduplicated devices sorted by IP.
func DiscoverWithContext(ctx context.Context) ([]Device, error) {
	conn, err := listenUDP(ICSPPort)
	if err != nil {
		return nil, err
	}

	defer func() { _ = conn.Close() }()

	// Poll the deadline every 200 ms so the context cancellation is noticed
	// promptly without blocking indefinitely on ReadFromUDP.
	const pollInterval = 200 * time.Millisecond

	return collect(conn, func() bool {
		select {
		case <-ctx.Done():
			return true
		default:
			_ = conn.SetReadDeadline(time.Now().Add(pollInterval))
			return false
		}
	}), nil
}

func listenUDP(port int) (*net.UDPConn, error) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port})
	if err != nil {
		return nil, fmt.Errorf("failed to listen on UDP port %d: %w", ICSPPort, err)
	}

	return conn, nil
}

// collect reads ICSP packets from conn until an error (e.g. deadline) occurs.
// If done is non-nil it is called before each read; returning true stops the
// loop immediately (used by DiscoverWithContext to check for cancellation).
func collect(conn *net.UDPConn, done func() bool) []Device {
	seen := make(map[string]Device)
	buf := make([]byte, 4096)

	for done == nil || !done() {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			if done != nil && isTimeout(err) {
				// Short deadline expired — check context and retry.
				continue
			}

			break
		}

		p, err := NewPacket(buf[:n])
		if err != nil {
			continue
		}

		device, err := p.Parse()
		if err != nil {
			continue
		}

		device.IP = remoteAddr.IP.String()
		// Deduplicate by MAC (last-seen wins).
		seen[device.MAC] = *device
	}

	result := make([]Device, 0, len(seen))
	for _, d := range seen {
		result = append(result, d)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].IP < result[j].IP
	})

	return result
}

// isTimeout reports whether err is a network timeout error.
func isTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
