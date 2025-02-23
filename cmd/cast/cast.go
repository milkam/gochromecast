package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/vjerci/gochromecast/pkg/chromecast"
	"github.com/vjerci/gochromecast/pkg/ip"
	"github.com/vjerci/gochromecast/pkg/mdns"
	"github.com/vjerci/gochromecast/pkg/server"
)

func main() {
	waitTime := flag.Int("mdns-seconds", 5, "specify number of seconds for searching the mdns (defaults to 5)")
	ipv6 := flag.Bool("ipv6", false, "should you search for devices using ipv6 instead of ipv4 (ipv4 is deafult)")
	targetDevice := flag.String("device", "", "specify a device name to cast to")

	flag.Parse()

	if len(*targetDevice) == 0 {
		panic("didnt specify a device name paramater which is required")
	}

	deviceToUse, err := getDevice(ipv6, waitTime, targetDevice)
	if err != nil {
		panic(fmt.Errorf("failed to find device %w", err))
	}

	ip, err := ip.GetLANIp()
	if err != nil {
		panic(fmt.Errorf("failed to get ip %s", ip))
	}

	log.Printf("resolving your lan ip  to %s", ip)

	const serverPort = ":8889"

	chromecastCtx, chromecastCancel := context.WithCancel(context.Background())
	chrmoecastClient := chromecast.New(chromecastCtx, &chromecast.Config{
		Device: deviceToUse,
	})
	defer chromecastCancel()

	go server.Start(serverPort)

	// wait for server to start as it needs to be online once we send to chromecast to start casting
	time.Sleep(1 * time.Second)

	err = chrmoecastClient.PlayMedia(chromecastCtx, chromecast.PlayMediaRequest{
		ChromeCastDeviceURI: deviceToUse.Url,
		MediaURL:            fmt.Sprintf("http://%s%s/files/playlist.m3u8", ip, serverPort),
		SubtitlesURL:        fmt.Sprintf("http://%s%s/files/subtitles.vtt", ip, serverPort),
	})
	if err != nil {
		panic(fmt.Errorf("failed to cast to tv %w", err))
	}

	time.Sleep(10 * time.Minute)
}

func getDevice(ipv6 *bool, waitTime *int, targetDevice *string) (mdns.Device, error) {
	mdnsCtx, mdnsCancel := context.WithCancel(context.Background())
	mdnsClient := mdns.New(mdnsCtx, &mdns.Config{
		IPv6: *ipv6,
	})

	mdnsClient.Start()

	time.Sleep(time.Duration(*waitTime) * time.Second)

	devicesChan := mdnsClient.GetDevices()

	devices := <-devicesChan

	mdnsCancel()

	for _, device := range devices {
		for _, name := range device.Names {
			if name == *targetDevice {
				return device, nil
			}
		}
	}

	return mdns.Device{}, fmt.Errorf("failed to find device for name '%s'", *targetDevice)
}
