package geography

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/shanth1/gotools/ops"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/domain"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type GeoProvider struct {
	geoDB *geoip2.Reader
}

func NewGeoProvider(cfg *config.Config) (ports.GeoProvider, error) {
	const op = "geography.NewGeoProvider"

	var db *geoip2.Reader
	var err error

	if cfg.GeoIP.Enabled {
		db, err = geoip2.Open(cfg.GeoIP.DBPath)
		if err != nil {
			return nil, ops.Wrap(op, ops.KindInternal, fmt.Errorf("open geoip2 database with path %q: %w", cfg.GeoIP.DBPath, err))
		}
	}

	return &GeoProvider{
		geoDB: db,
	}, nil
}

func (g *GeoProvider) Lookup(_ context.Context, ipStr string) (*domain.GeoLocation, error) {
	const op = "geography.GeoProvider.Lookup"

	if ipStr == "127.0.0.1" || ipStr == "::1" {
		return &domain.GeoLocation{
			Country: "Loopback",
			City:    "Localhost",
		}, nil
	}

	if g.geoDB == nil {
		return nil, ops.Wrap(op, ops.KindInternal, errors.New("geo db is nil"))
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, ops.Wrap(op, ops.KindInvalid, fmt.Errorf("parse ip: %q", ipStr))
	}

	if isPrivateIP(ip) {
		return &domain.GeoLocation{
			Country: "Private",
			City:    "LAN",
		}, nil
	}

	record, err := g.geoDB.City(ip)
	if err != nil {
		return nil, ops.Wrap(op, ops.KindInternal, fmt.Errorf("get geoip2 record with ip %q: %w", ip.String(), err))
	}

	country := record.Country.IsoCode
	city := record.City.Names["en"]

	if country == "" {
		country = "XX"
	}
	if city == "" && country != "XX" {
		city = "Unknown"
	}

	return &domain.GeoLocation{
		Country: country,
		City:    city,
	}, nil
}

func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return true
	}
	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}
	return false
}
