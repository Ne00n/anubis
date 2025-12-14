package thoth

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// LocalDBs is a client that looks up IP addresses in local GeoIP and ASN databases.
type LocalDBs struct {
	geoDB *geoip2.Reader
	asnDB *geoip2.Reader
}

// NewLocalDBs creates a new LocalDBs client.
func NewLocalDBs(geoPath, asnPath string) (*LocalDBs, error) {
	var geoDB, asnDB *geoip2.Reader
	var err error

	if geoPath != "" {
		geoDB, err = geoip2.Open(geoPath)
		if err != nil {
			return nil, fmt.Errorf("could not open geoip database at %s: %w", geoPath, err)
		}
	}

	if asnPath != "" {
		asnDB, err = geoip2.Open(asnPath)
		if err != nil {
			return nil, fmt.Errorf("could not open asn database at %s: %w", asnPath, err)
		}
	}

	return &LocalDBs{geoDB: geoDB, asnDB: asnDB}, nil
}

// Lookup looks up an IP address in the local GeoIP and ASN databases.
func (l *LocalDBs) Lookup(ipStr string) (*GeoIPInfo, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	info := &GeoIPInfo{
		Announced: true, // Local database always has announced IPs
	}

	if l.geoDB != nil {
		record, err := l.geoDB.Country(ip)
		if err == nil && record != nil {
			info.CountryCode = record.Country.IsoCode
		}
	}

	if l.asnDB != nil {
		record, err := l.asnDB.ASN(ip)
		if err == nil && record != nil {
			info.ASNumber = uint32(record.AutonomousSystemNumber)
		}
	}

	return info, nil
}

// Close closes the underlying databases.
func (l *LocalDBs) Close() error {
	var err error
	if l.geoDB != nil {
		err = l.geoDB.Close()
	}
	if l.asnDB != nil {
		if err != nil {
			// Don't overwrite the first error
			l.asnDB.Close()
		} else {
			err = l.asnDB.Close()
		}
	}
	return err
}

// GeoIPInfo is a mock of the iptoasnv1.LookupResponse object.
type GeoIPInfo struct {
	CountryCode string
	ASNumber    uint32
	Announced   bool
}

// GetCountryCode returns the country code.
func (i *GeoIPInfo) GetCountryCode() string {
	return i.CountryCode
}

// GetASNumber returns the ASN number.
func (i *GeoIPInfo) GetASNumber() uint32 {
	return i.ASNumber
}

// GetAnnounced returns whether the IP is announced.
func (i *GeoIPInfo) GetAnnounced() bool {
	return i.Announced
}
