package mdns

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/mdns"
)

var DefaultChromecastPort = "8009"

type Client struct {
	config *Config
	ctx    context.Context

	newDevice        chan *mdns.ServiceEntry
	getDevicesC      chan chan []Device
	availableDevices []Device
}

type Config struct {
	IPv6 bool
}

func New(ctx context.Context, config *Config) *Client {
	return &Client{
		config: config,
		ctx:    ctx,

		getDevicesC: make(chan chan []Device, 100),
		newDevice:   make(chan *mdns.ServiceEntry, 100),
	}
}

func (client *Client) Start() {
	go client.listen()
	go client.startDiscovery()
}

func (client *Client) listen() {
	for {
		select {
		case <-client.ctx.Done():
			return
		case serviceEntry := <-client.newDevice:
			if serviceEntry == nil {
				log.Printf("mdns got device which is nil")
				continue
			}

			if !client.hasAddress(serviceEntry) {
				continue
			}

			log.Printf("mdns got device %s %s", serviceEntry.Name, fmt.Sprintf("%s:%d", serviceEntry.AddrV4.String(), serviceEntry.Port))

			client.addDevice(serviceEntry)
		case responseChan := <-client.getDevicesC:
			responseChan <- client.currentDevices()
		}
	}
}

type Device struct {
	Names []string
	Url   string
}

func (client *Client) GetDevices() chan []Device {
	responseChan := make(chan []Device, 1)

	go func() {
		client.getDevicesC <- responseChan
	}()

	return responseChan
}

func (client *Client) currentDevices() []Device {
	returnedDevices := []Device{}

	for _, device := range client.availableDevices {
		returnedDevices = append(returnedDevices, Device{
			Names: append([]string{}, device.Names...),
			// there is a service.Port but it returns wrong port such as 6466s
			Url: fmt.Sprintf("%s:%s", device.Url, DefaultChromecastPort),
		})
	}

	return returnedDevices
}

func (client *Client) addDevice(mdnsEntry *mdns.ServiceEntry) {
	log.Printf("trying to add device %s", mdnsEntry.Name)

	deviceExists := false

	name := client.toHumanDeviceName(mdnsEntry.Name)

	for i, v := range client.availableDevices {
		if v.Url == client.getAddress(mdnsEntry) {
			deviceExists = true

			log.Printf("device already exists with different name, adding name %s", name)

			client.availableDevices[i].Names = append(client.availableDevices[i].Names, client.toHumanDeviceName(mdnsEntry.Name))
			break
		}
	}

	if !deviceExists {
		log.Printf("device doesn't exist adding it to devices list name is '%s'", name)

		client.availableDevices = append(client.availableDevices, Device{
			Names: []string{name},
			Url:   client.getAddress(mdnsEntry),
		})
	}
}

func (client *Client) toHumanDeviceName(name string) string {
	if strings.HasSuffix(name, "._androidtvremote2._tcp.local.") {
		name, _ = strings.CutSuffix(name, "._androidtvremote2._tcp.local.")
	}

	if index := strings.Index(name, "."); index != -1 {
		name = name[:index]
	}

	return name
}

func (client *Client) getAddress(serviceEntry *mdns.ServiceEntry) string {
	if client.config.IPv6 {
		return serviceEntry.AddrV6.String()
	} else {
		return serviceEntry.AddrV4.String()
	}
}

func (client *Client) hasAddress(serviceEntry *mdns.ServiceEntry) bool {
	if client.config.IPv6 {
		if len(serviceEntry.AddrV6) == 0 {
			log.Printf("mdns skipping device due to no ipv6 address %#v", serviceEntry)
			return false
		}
	} else {
		if len(serviceEntry.AddrV4) == 0 {
			log.Printf("mdns skipping device due to no ipv4 address %#v", serviceEntry)
			return false
		}
	}

	return true
}
