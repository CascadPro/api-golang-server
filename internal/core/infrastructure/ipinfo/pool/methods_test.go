package core_ipinfo_pool

import (
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ipinfo/go/v2/ipinfo"
)

type fakeIPInfoClient struct {
	calls atomic.Int32
}

func (f *fakeIPInfoClient) GetIPInfo(net.IP) (*ipinfo.Core, error) {
	f.calls.Add(1)
	time.Sleep(25 * time.Millisecond)
	return &ipinfo.Core{}, nil
}

func TestGetIPInfoCoalescesConcurrentRequestsForSameIP(t *testing.T) {
	client := &fakeIPInfoClient{}
	pool := &ConnectionPool{
		client:  client,
		timeout: time.Second,
		flights: make(map[string]*ipInfoCall),
	}

	const workers = 12
	var wg sync.WaitGroup
	errs := make(chan error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			info, err := pool.GetIPInfo(net.ParseIP("8.8.8.8"))
			if err == nil && info == nil {
				errs <- net.InvalidAddrError("nil result")
				return
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := client.calls.Load(); got != 1 {
		t.Fatalf("upstream called %d times, want 1", got)
	}
}
