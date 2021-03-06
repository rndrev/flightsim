package simulator

import (
	"context"
	"fmt"
	"math/rand"
	"net"
)

var (
	scanPorts = []int{21, 22, 23, 25, 80, 88, 111, 135, 139, 143, 389, 443, 445, 1433, 1521, 3306, 3389, 5432, 5900, 6000, 8443}

	// IP ranges are for TEST-NET-[123] networks, as per RFC 5737.
	// These IPs should be safe to scan as they're not assigned.
	scanIPRanges = []*net.IPNet{
		{
			IP:   net.IPv4(192, 0, 2, 0),
			Mask: net.CIDRMask(24, 32),
		},
		{
			IP:   net.IPv4(198, 51, 100, 0),
			Mask: net.CIDRMask(24, 32),
		},
		{
			IP:   net.IPv4(203, 0, 113, 0),
			Mask: net.CIDRMask(24, 32),
		},
	}
)

func randIP(network *net.IPNet) net.IP {
	randIP := make(net.IP, len(network.Mask))
	rand.Read(randIP)

	// reverse mask and map randIP, so it only contains bits
	// which can be randomized
	mask := make(net.IPMask, len(network.Mask))
	for n := range mask {
		mask[n] = ^network.Mask[n]
	}
	randIP = randIP.Mask(mask)

	netIP := network.IP.To16()[16-len(randIP):]
	for n := range randIP {
		randIP[n] = randIP[n] | netIP[n]
	}

	return randIP
}

// PortScan simulator.
type PortScan struct {
}

// NewPortScan creates port scan simulator.
func NewPortScan() *PortScan {
	return &PortScan{}
}

// Simulate port scanning for given host.
func (*PortScan) Simulate(ctx context.Context, extIP net.IP, host string) error {
	d := &net.Dialer{
		LocalAddr: &net.TCPAddr{IP: extIP},
	}

	conn, err := d.DialContext(ctx, "tcp", host)
	if err != nil {
		return err
	}
	conn.Close()

	return nil
}

// Hosts returns host:port generated from RFC 5737 addresses.
func (s *PortScan) Hosts(size int) ([]string, error) {
	var hosts []string

	// for each network generate size IPs and add all ports;
	// total number of hosts will be up to size * networks * ports
	for _, network := range scanIPRanges {
		// TODO: make caller responsible for deduplication
		dedup := make(map[string]bool)

		for k := 0; k < size; k++ {
			ip := randIP(network)

			key := ip.String()
			if dedup[key] {
				continue
			}
			dedup[key] = true

			for _, port := range scanPorts {
				hosts = append(hosts, fmt.Sprintf("%s:%d", ip, port))
			}
		}
	}

	return hosts, nil
}
