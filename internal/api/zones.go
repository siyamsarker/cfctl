package api

import (
	"context"
	"fmt"

	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

// ListZones retrieves all zones for the account
func (c *Client) ListZones(ctx context.Context) ([]cloudflare.Zone, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	zones, err := c.api.ListZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("list zones: %w", err)
	}

	result := make([]cloudflare.Zone, len(zones))
	for i, z := range zones {
		result[i] = cloudflare.Zone{
			ID:     z.ID,
			Name:   z.Name,
			Status: z.Status,
			Plan: cloudflare.Plan{
				Name: z.Plan.Name,
			},
		}
	}

	return result, nil
}

// GetZone retrieves a specific zone by ID
func (c *Client) GetZone(ctx context.Context, zoneID string) (*cloudflare.Zone, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	z, err := c.api.ZoneDetails(ctx, zoneID)
	if err != nil {
		return nil, fmt.Errorf("get zone: %w", err)
	}

	return &cloudflare.Zone{
		ID:     z.ID,
		Name:   z.Name,
		Status: z.Status,
		Plan: cloudflare.Plan{
			Name: z.Plan.Name,
		},
	}, nil
}
