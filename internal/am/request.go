package am

import (
	"net"
	"net/http"
	"strings"
)

func GetIPAddress(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		split := strings.Split(xForwardedFor, ",")
		if len(split) > 0 {
			return strings.TrimSpace(split[0])
		}
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// GetGeoData retrieves geographical information based on an IP address.
func GetGeoData(r *http.Request) GeoPoint {
	// Placeholder implementation
	// TODO: Research a good geoip library
	return GeoPoint{
		Lat: -6.1751,
		Lng: 106.6297,
	}
}
