package chromecast

import "github.com/vjerci/gochromecast/pkg/chromecast/proto/castchannel"

var NamespaceConnection = "urn:x-cast:com.google.cast.tp.connection"

type PayloadConnect struct {
	Type       string                    `json:"type"`
	Origin     map[string]string         `json:"origin"`
	UserAgent  string                    `json:"userAgent"`
	SenderInfo *PayloadConnectSenderInfo `json:"senderInfo"`
}

type PayloadConnectSenderInfo struct {
	SDKType        int    `json:"sdkType"`
	Version        string `json:"version"`
	BrowserVersion string `json:"browserVersion"`
	Platform       int    `json:"platform"`
	SystemVersion  string `json:"systemVersion"`
	ConnectionType int    `json:"connectionType"`
}

var payloadConnect = PayloadConnect{
	Type:      "CONNECT",
	UserAgent: "GoChromecast",
	Origin:    make(map[string]string),
	SenderInfo: &PayloadConnectSenderInfo{
		SDKType:        2,
		Version:        "15.605.1.3",
		BrowserVersion: "44.0.2403.30",
		Platform:       4,
		ConnectionType: 1,
		SystemVersion:  "Macintosh; Intel Mac OS X10_10_3",
	},
}

func (client *Client) sendConnectToDevice(sender *Sender) {
	sender.SendMsg(SenderMessage{
		Proto: &castchannel.CastMessage{
			ProtocolVersion: ProtocolVersion,
			SourceId:        &SenderID,
			DestinationId:   &ReceiverID,
			Namespace:       &NamespaceConnection,
			PayloadType:     PayloadTypeString,
		},
		JsonData: payloadConnect,
	})
}
