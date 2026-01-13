package api

import (
	"context"
	"fmt"

	cf "github.com/cloudflare/cloudflare-go"
	"github.com/siyamsarker/cfctl/pkg/cloudflare"
)

// PurgeCache purges cache based on the provided request
func (c *Client) PurgeCache(ctx context.Context, zoneID string, req cloudflare.PurgeRequest) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	// Build the purge request based on the type
	var purgeReq cf.PurgeCacheRequest

	if req.PurgeEverything {
		// Purge everything
		purgeReq = cf.PurgeCacheRequest{
			Everything: true,
		}
	} else if len(req.Files) > 0 {
		// Purge by files/URLs
		purgeReq = cf.PurgeCacheRequest{
			Files: req.Files,
		}
	} else if len(req.Hosts) > 0 {
		// Purge by hosts
		purgeReq = cf.PurgeCacheRequest{
			Hosts: req.Hosts,
		}
	} else if len(req.Tags) > 0 {
		// Purge by tags (Enterprise only)
		purgeReq = cf.PurgeCacheRequest{
			Tags: req.Tags,
		}
	} else if len(req.Prefixes) > 0 {
		// Purge by prefixes
		purgeReq = cf.PurgeCacheRequest{
			Prefixes: req.Prefixes,
		}
	} else {
		return fmt.Errorf("no purge parameters provided")
	}

	_, err := c.api.PurgeCache(ctx, zoneID, purgeReq)
	if err != nil {
		return fmt.Errorf("purge cache: %w", err)
	}

	return nil
}
