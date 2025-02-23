package chromecast

import "github.com/vjerci/gochromecast/pkg/chromecast/proto/castchannel"

var PayloadTypeLaunch = "LAUNCH"

type PayloadLaunch struct {
	Type      string `json:"type"`
	AppID     string `json:"appId"`
	RequestID int    `json:"requestId"`
}

func (client *Client) sendLaunchRecieverAppMsg(sender *Sender, requestID int) {
	sender.SendMsg(SenderMessage{
		Proto: &castchannel.CastMessage{
			ProtocolVersion: ProtocolVersion,
			SourceId:        &SenderID,
			DestinationId:   &ReceiverID,
			Namespace:       &NamespaceReceiver,
			PayloadType:     PayloadTypeString,
		},
		JsonData: PayloadLaunch{
			Type:      PayloadTypeLaunch,
			AppID:     DefaultMediaAppID,
			RequestID: requestID,
		},
	})
}
