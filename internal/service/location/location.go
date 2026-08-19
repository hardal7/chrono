package location

import (
	"net/netip"

	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/oschwald/geoip2-golang/v2"
)

func IPToCountry(ip string) string {
	logger.Info("Linking IP address to country")

	db, err := geoip2.Open("/srv/GeoLite2-City.mmdb")
	if err != nil {
		logger.Error("Failed to load geolocation database", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			logger.Warn("Failed to close location database")
		}
	}()

	address, err := netip.ParseAddr(ip)
	if err != nil {
		logger.Error("Failed to parse IP address", err)
	}
	record, err := db.Country(address)
	if err != nil {
		logger.Error("Failed to find location of IP address", err)
	}
	if !record.HasData() {
		logger.Error("No data found for this IP")
		return ""
	}

	country := record.Country.Names.English
	logger.Info("Linked IP address to country", "country", country)
	return country
}
