package core_ipinfo_pool

import (
	"fmt"
	"net"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/ipinfo/go/v2/ipinfo"
)

type ipInfoCall struct {
	done chan struct{}
	info *ipinfo.Core
	err  error
}

// GetIPInfo coalesces concurrent lookups for the same IP within this process.
// This complements Redis caching: when many requests miss the cache at once,
// only one request per IP reaches the upstream IPInfo client.
func (p *ConnectionPool) GetIPInfo(ip net.IP) (*ipinfo.Core, error) {
	if ip == nil {
		return nil, fmt.Errorf("%v: %w", net.InvalidAddrError("IP address is nil"), core_errors.ErrInvalidArgument)
	}

	key := ip.String()
	p.flightsMu.Lock()

	if active, ok := p.flights[key]; ok {
		p.flightsMu.Unlock()

		<-active.done

		return active.info, active.err
	}

	active := &ipInfoCall{done: make(chan struct{})}
	if p.flights == nil {
		p.flights = make(map[string]*ipInfoCall)
	}

	p.flights[key] = active

	p.flightsMu.Unlock()

	active.info, active.err = p.client.GetIPInfo(ip)

	p.flightsMu.Lock()

	delete(p.flights, key)

	close(active.done)

	p.flightsMu.Unlock()

	return active.info, active.err
}
