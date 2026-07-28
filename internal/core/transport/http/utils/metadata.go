package core_http_utils

import (
	"context"
	"fmt"
	"net"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_ipinfo_client "github.com/CascadePro/api-golang-server/internal/core/repository/ipinfo/client"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
)

const unknown = "Unknown"

func GetSessionMetadata(
	ctx context.Context,
	ipClient core_ipinfo_client.RepositoryMethods,
	ip net.IP,
	userAgent *core_http_request.UserAgent,
) (domain.SessionMetadata, error) {
	cfg, err := core_config.NewConfig()
	if err != nil {
		return domain.SessionMetadata{}, fmt.Errorf("get core config: %w", err)
	}

	var location domain.SessionMetadataLocation
	if cfg.Connection == core_config.ConnectionOffline {
		location = domain.NewSessionMetadataLocation(unknown, unknown, 0, 0)
	} else {
		ipinfo, err := ipClient.Lookup(ctx, ip)
		if err != nil {
			return domain.SessionMetadata{}, fmt.Errorf("get ipinfo from repository: %w", err)
		}

		location = domain.NewSessionMetadataLocation(ipinfo.Country, ipinfo.City, ipinfo.Lat, ipinfo.Lng)
	}

	device := domain.NewSessionMetadataDevice(
		userAgent.OS,
		userAgent.Model,
		userAgent.AppName,
		userAgent.Type,
		userAgent.AppVersion,
	)

	return domain.NewSessionMetadata(location, device), nil
}
