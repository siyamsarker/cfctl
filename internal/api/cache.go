package api

import (
	"context"
	"fmt"

	cfv6 "github.com/cloudflare/cloudflare-go/v6"
	"github.com/cloudflare/cloudflare-go/v6/cache"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

// PurgeCache purges cache based on the provided request
func (c *Client) PurgeCache(ctx context.Context, zoneID string, req cloudflare.PurgeRequest) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build the purge request based on the type
	var body cache.CachePurgeParamsBody

	if req.PurgeEverything {
		// Purge everything
		body = cache.CachePurgeParamsBody{
			PurgeEverything: cfv6.F(true),
		}
	} else if len(req.Files) > 0 {
		// Purge by files/URLs
		body = cache.CachePurgeParamsBody{
			Files: cfv6.F[interface{}](req.Files),
		}
	} else if len(req.Hosts) > 0 {
		// Purge by hosts
		body = cache.CachePurgeParamsBody{
			Hosts: cfv6.F[interface{}](req.Hosts),
		}
	} else if len(req.Tags) > 0 {
		// Purge by tags (Enterprise only)
		body = cache.CachePurgeParamsBody{
			Tags: cfv6.F[interface{}](req.Tags),
		}
	} else if len(req.Prefixes) > 0 {
		// Purge by prefixes
		body = cache.CachePurgeParamsBody{
			Prefixes: cfv6.F[interface{}](req.Prefixes),
		}
	} else {
		return fmt.Errorf("no purge parameters provided")
	}

	purgeParams := cache.CachePurgeParams{
		ZoneID: cfv6.F(zoneID),
		Body:   body,
	}

	_, err := c.api.Cache.Purge(ctx, purgeParams)
	if err != nil {
		return fmt.Errorf("purge cache: %w", err)
	}

	return nil
}
