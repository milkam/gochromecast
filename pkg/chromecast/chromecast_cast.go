package chromecast

import (
	"context"
	"crypto/tls"
	"errors"
	"log"

	"github.com/vjerci/gochromecast/pkg/chromecast/proto/castchannel"
)

// chromecast general ids
var SenderID = "sender-vjerci"
var ReceiverID = "receiver-0"      //hardcoded needs to be this value
var DefaultMediaAppID = "CC1AD845" // hardcoded default media reciever

var ProtocolVersion = castchannel.CastMessage_CASTV2_1_0.Enum()

var PayloadTypeBinary = castchannel.CastMessage_BINARY.Enum()
var PayloadTypeString = castchannel.CastMessage_STRING.Enum()

var ErrPlayMediaDial = errors.New("dialing failed")

var ErrPlayMediaTransportID = errors.New("getting transport id failed")

type PlayMediaRequest struct {
	ChromeCastDeviceURI string // 192.168.0.102:8009
	MediaURL            string // full media url "http://192.168.8.132:1113/livestream.m3u8"
	SubtitlesURL        string // full subtitles url "http://192.168.8.132:1113/sub_df39804e-65e8-498e-8796-e45950e75ffd.vtt"
}

func (client *Client) PlayMedia(ctx context.Context, req PlayMediaRequest) error {
	conn, err := tls.Dial("tcp", req.ChromeCastDeviceURI, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return errors.Join(ErrPlayMediaDial, err)
	}

	tlsReader := NewReader(ctx, conn)
	tlsReader.Start()

	sender := NewSender(ctx, conn)
	sender.Start()

	requestCounter := NewRequestCounter(ctx)
	requestCounter.Start()

	// start ping pong
	client.startPingHandler(ctx, sender, requestCounter, tlsReader)

	// send connect message
	log.Printf("sending connect to device")
	client.sendConnectToDevice(sender)

	requestIDLaunchRecieverApp := <-requestCounter.GetRequestCounter()

	statusSub := client.SubscribeRecieverStatus(ctx, tlsReader, requestIDLaunchRecieverApp)

	log.Printf("sending launch reciever app msg")
	// send launch default reciever app msg
	client.sendLaunchRecieverAppMsg(sender, requestIDLaunchRecieverApp)

	// wait for status to show which contains our transport id, or simply time out
	transportID, err := client.WaitForTransportID(statusSub)
	if err != nil {
		return errors.Join(ErrPlayMediaTransportID, err)
	}

	// send connect to transport id/reciever app
	log.Printf("sending connect to reciver app using transport ID %s", transportID)
	client.SendConnectToApp(sender, transportID)

	log.Printf("sending load media")
	// send loadMedia request
	client.SendLoadMedia(requestCounter, req, transportID, sender)

	return nil
}
