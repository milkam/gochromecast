package chromecast

import (
	"context"

	"github.com/google/uuid"
)

func (client *Client) startPingHandler(ctx context.Context, sender *Sender, requestCounter *RequestCounter, reader *Reader) {
	pingHandler := &PingHandler{
		ctx:           ctx,
		id:            uuid.NewString(),
		sender:        sender,
		counterGetter: requestCounter,
	}
	reader.Subscribe(pingHandler)

	pingHandler.Start()
}
