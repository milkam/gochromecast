package chromecast

import (
	"context"

	"github.com/milkam/gochromecast/pkg/mdns"
)

type Client struct {
	config *Config
	ctx    context.Context
}

type Config struct {
	Device mdns.Device
}

func New(ctx context.Context, config *Config) *Client {
	return &Client{
		config: config,
		ctx:    ctx,
	}
}
