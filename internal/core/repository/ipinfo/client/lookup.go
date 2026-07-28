package core_ipinfo_client

import (
	"context"
	"fmt"
	"net"
	"strings"

	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
)

func (r *Repository) Lookup(ctx context.Context, ip net.IP) (InfoModel, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	info, err := r.pool.GetIPInfo(ip.To4())
	if err != nil {
		return InfoModel{}, fmt.Errorf("ipinfo lookup: %w", err)
	}

	model := InfoModel{
		IP:       ip,
		City:     info.City,
		Country:  info.CountryName,
		Region:   info.Region,
		Timezone: info.Timezone,
	}

	if len(info.Location) > 0 {
		coords := strings.Split(info.Location, ",")
		if len(coords) == 2 {
			model.Lat = core_utils.ParseFloat(coords[0])
			model.Lng = core_utils.ParseFloat(coords[1])
		}
	}

	return model, nil
}
