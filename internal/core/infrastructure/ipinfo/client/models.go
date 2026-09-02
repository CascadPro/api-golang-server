package core_ipinfo_client

import "net"

type InfoModel struct {
	IP       net.IP  `json:"ip"`
	City     string  `json:"city,omitempty"`
	Region   string  `json:"region,omitempty"`
	Country  string  `json:"country,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Timezone string  `json:"timezone,omitempty"`
}
