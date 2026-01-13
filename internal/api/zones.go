package api

import (
	"context"
	"fmt"

	cfv6 "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/zones"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

// ListZones retrieves all zones for the account
func (c *Client) ListZones(ctx context.Context) ([]cloudflare.Zone, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// List zones using v6 API with auto-pagination
	var allZones []cloudflare.Zone

	pager := c.api.Zones.ListAutoPaging(ctx, zones.ZoneListParams{
		PerPage: cfv6.F(float64(50)),
	})

	for pager.Next() {
		z := pager.Current()
		allZones = append(allZones, cloudflare.Zone{
			ID:     z.ID,
			Name:   z.Name,
			Status: string(z.Status),
			Plan: cloudflare.Plan{
				Name: z.Plan.Name,
			},
		})
	}

	if err := pager.Err(); err != nil {
		return nil, fmt.Errorf("list zones: %w", err)
	}

	return allZones, nil
}

// GetZone retrieves a specific zone by ID
func (c *Client) GetZone(ctx context.Context, zoneID string) (*cloudflare.Zone, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	z, err := c.api.Zones.Get(ctx, zones.ZoneGetParams{
		ZoneID: cfv6.F(zoneID),
	})
	if err != nil {
		return nil, fmt.Errorf("get zone: %w", err)
	}

	return &cloudflare.Zone{
		ID:     z.ID,
		Name:   z.Name,
		Status: string(z.Status),
		Plan: cloudflare.Plan{
			Name: z.Plan.Name,
		},
	}, nil
}
