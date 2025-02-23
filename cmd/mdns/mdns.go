package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/vjerci/gochromecast/pkg/mdns"
)

func main() {
	waitTime := flag.Int("mdns-seconds", 5, "specify number of seconds for searching the mdns (defaults to 5)")
	ipv6 := flag.Bool("ipv6", false, "should you search for devices using ipv6 instead of ipv4 (ipv4 is deafult)")

	ctx, ctxCancel := context.WithCancel(context.Background())
	mdns := mdns.New(ctx, &mdns.Config{
		IPv6: *ipv6,
	})

	mdns.Start()

	time.Sleep(time.Duration(*waitTime) * time.Second)

	devicesChan := mdns.GetDevices()

	devices := <-devicesChan

	log.Printf("closing program")

	for _, device := range devices {
		log.Printf("got devices on your network names:'%#v' url:'%s'", device.Names, device.Url)
	}

	ctxCancel()
}
