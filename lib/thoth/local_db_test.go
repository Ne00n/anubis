package thoth

import (
	"net/http/httptest"
	"os"
	"testing"
)

const (
	geoLite2CountryDbPath = "../../test/testdata/GeoLite2-Country.mmdb"
	geoLite2AsnDbPath     = "../../test/testdata/GeoLite2-ASN.mmdb"
)

// To run this test, you need to download the GeoLite2 Country and ASN databases from MaxMind
// and place them in test/testdata/
// You can download them from https://dev.maxmind.com/geoip/geolite2-free-geolocation-data

func TestLocalGeoIPChecker(t *testing.T) {
	if _, err := os.Stat(geoLite2CountryDbPath); os.IsNotExist(err) {
		t.Skipf("skipping test; GeoLite2 Country database not found at %s", geoLite2CountryDbPath)
	}

	client, err := NewLocal(geoLite2CountryDbPath, "")
	if err != nil {
		t.Fatalf("failed to create local thoth client: %v", err)
	}
	defer client.Close()

	checker := client.GeoIPCheckerFor([]string{"us"}).(*GeoIPChecker)

	cases := []struct {
		name      string
		ip        string
		wantMatch bool
	}{
		{"Cloudflare US", "1.1.1.1", true},
		{"Google US", "8.8.8.8", true},
		{"OpenDNS US", "208.67.222.222", true},
		{"Japanese DNS", "202.248.37.74", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Real-Ip", tc.ip)

			match, err := checker.Check(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if match != tc.wantMatch {
				t.Errorf("want match %v, got %v", tc.wantMatch, match)
			}
		})
	}
}

func TestLocalASNChecker(t *testing.T) {
	if _, err := os.Stat(geoLite2AsnDbPath); os.IsNotExist(err) {
		t.Skipf("skipping test; GeoLite2 ASN database not found at %s", geoLite2AsnDbPath)
	}

	client, err := NewLocal("", geoLite2AsnDbPath)
	if err != nil {
		t.Fatalf("failed to create local thoth client: %v", err)
	}
	defer client.Close()

	checker := client.ASNCheckerFor([]uint32{13335}).(*ASNChecker)

	cases := []struct {
		name      string
		ip        string
		wantMatch bool
	}{
		{"Cloudflare US", "1.1.1.1", true},
		{"Google US", "8.8.8.8", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("X-Real-Ip", tc.ip)

			match, err := checker.Check(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if match != tc.wantMatch {
				t.Errorf("want match %v, got %v", tc.wantMatch, match)
			}
		})
	}
}
