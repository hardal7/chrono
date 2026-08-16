package location

import (
	"net/netip"

	"github.com/hardal7/chrono/internal/util/logger"
	"github.com/oschwald/geoip2-golang/v2"
)

func IPToLocation(ip string) string {
	logger.Info("Linking IP address to geolocation")

	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		logger.Error("Failed to load geolocation database", err)
	}
	defer db.Close()

	address, err := netip.ParseAddr(ip)
	if err != nil {
		logger.Error("Failed to parse IP address", err)
	}
	record, err := db.City(address)
	if err != nil {
		logger.Error("Failed to find location of IP address", err)
	}
	if !record.HasData() {
		logger.Error("No data found for this IP")
		return ""
	}

	city := record.City.Names.English
	logger.Info("Linked IP address to geolocation", "city", city)
	return city
}
