package chromecast

import "github.com/vjerci/gochromecast/pkg/chromecast/proto/castchannel"

func (client *Client) SendConnectToApp(sender *Sender, transportID string) {
	sender.SendMsg(SenderMessage{
		Proto: &castchannel.CastMessage{
			ProtocolVersion: ProtocolVersion,
			SourceId:        &SenderID,
			DestinationId:   &transportID,
			Namespace:       &NamespaceConnection,
			PayloadType:     PayloadTypeString,
		},
		JsonData: payloadConnect,
	})
}
