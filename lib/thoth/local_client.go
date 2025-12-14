package thoth

import (
	"context"
	"fmt"

	iptoasnv1 "github.com/TecharoHQ/thoth-proto/gen/techaro/thoth/iptoasn/v1"
	"google.golang.org/grpc"
)

// localIpToASNClient is a wrapper around LocalDBs that implements the iptoasnv1.IpToASNServiceClient interface.
type localIpToASNClient struct {
	localDBs *LocalDBs
}

// newLocalIpToASNClient creates a new localIpToASNClient.
func newLocalIpToASNClient(geoPath, asnPath string) (*localIpToASNClient, error) {
	localDBs, err := NewLocalDBs(geoPath, asnPath)
	if err != nil {
		return nil, fmt.Errorf("could not create local geoip client: %w", err)
	}
	return &localIpToASNClient{localDBs: localDBs}, nil
}

// Lookup looks up an IP address in the local GeoIP and ASN databases.
func (c *localIpToASNClient) Lookup(ctx context.Context, in *iptoasnv1.LookupRequest, opts ...grpc.CallOption) (*iptoasnv1.LookupResponse, error) {
	info, err := c.localDBs.Lookup(in.GetIpAddress())
	if err != nil {
		return nil, err
	}

	return &iptoasnv1.LookupResponse{
		CountryCode: info.GetCountryCode(),
		AsNumber:    info.GetASNumber(),
		Announced:   info.GetAnnounced(),
	}, nil
}

func (c *localIpToASNClient) Close() error {
	return c.localDBs.Close()
}
