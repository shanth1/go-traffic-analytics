package geography

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/shanth1/gotrace/internal/config"
	"github.com/shanth1/gotrace/internal/core/ports"
)

type GeoIPRepo struct {
	geoDB *geoip2.Reader
}

func NewGeoIPRepo(cfg *config.Config) (ports.GeoIPRepository, error) {
	var db *geoip2.Reader
	var err error

	if cfg.GeoIP.Enabled {
		db, err = geoip2.Open(cfg.GeoIP.DBPath)
		if err != nil {
			return nil, fmt.Errorf("open geoip2 database: %w", err)
		}
	}

	return &GeoIPRepo{
		geoDB: db,
	}, nil
}

func (g *GeoIPRepo) GetInfo(_ context.Context, ipStr string) (country, city string, err error) {
	if ipStr == "127.0.0.1" || ipStr == "::1" {
		return "Local", "Host", nil
	}

	if g.geoDB == nil {
		return "", "", errors.New("geo db is nil")
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", fmt.Errorf("parse ip: %s", ipStr)
	}

	record, err := g.geoDB.City(ip)
	if err != nil {
		return "", "", fmt.Errorf("get geoip2 record: %w", err)
	}

	country = record.Country.IsoCode
	city = record.City.Names["en"]

	return country, city, nil
}
