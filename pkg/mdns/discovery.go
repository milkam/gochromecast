package mdns

import (
	"errors"
	"log"
	"time"

	"github.com/hashicorp/mdns"
)

var ErrStartDiscoveryBrowse = errors.New("browsing zeroconf failed")

var DiscoveryServices = []string{"_googlecast._tcp", "_androidtvremote2._tcp"}

func (client *Client) startDiscovery() error {
	for _, service := range DiscoveryServices {
		localService := service

		localServiceReceive := make(chan *mdns.ServiceEntry, 10)
		go func() {
			for {
				select {
				case <-client.ctx.Done():
					return
				case service := <-localServiceReceive:
					client.newDevice <- service
				}
			}
		}()

		go func() {
			log.Printf("mdns listening for %s", localService)

			queryConfig := &mdns.QueryParam{
				Service:     localService,
				Entries:     localServiceReceive,
				Domain:      "local",
				DisableIPv6: true,
				// 10 years of timeout
				Timeout: time.Hour * 24 * 31 * 12 * 10,
			}

			if client.config.IPv6 {
				queryConfig.DisableIPv6 = false
			}

			err := mdns.Query(queryConfig)
			if err != nil {
				log.Printf("browsing error for service '%s' '%s'", localService, errors.Join(ErrStartDiscoveryBrowse, err))
			}
		}()
	}

	return nil
}
