package core_ipinfo_pool

import (
	"net"
	"time"

	"github.com/ipinfo/go/v2/ipinfo"
)

func (p *ConnectionPool) GetIPInfo(ip net.IP) (*ipinfo.Core, error) {
	return p.client.GetIPInfo(ip)
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}
