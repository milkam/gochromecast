package chromecast

import (
	"context"
	"time"

	"github.com/milkam/gochromecast/pkg/chromecast/proto/castchannel"
)

var NamespaceHeartbeat = "urn:x-cast:com.google.cast.tp.heartbeat"

var PayloadTypePing = "PING"
var PayloadTypePong = "PONG"

type CastMessageSender interface {
	SendMsg(SenderMessage)
}

type CounterGetter interface {
	GetRequestCounter() chan int
}

type PingHandler struct {
	ctx           context.Context
	id            string
	sender        CastMessageSender
	counterGetter CounterGetter
}

func (pingHandler *PingHandler) Start() {
	go pingHandler.startPinging()
}

func (pingHandler *PingHandler) startPinging() {
	go pingHandler.sendPing()

	t := time.NewTicker(4 * time.Second)
	for {
		select {
		case <-pingHandler.ctx.Done():
			return
		case <-t.C:
			pingHandler.sendPing()
		}
	}
}

func (pingHandler *PingHandler) sendPing() {
	requestIDChan := pingHandler.counterGetter.GetRequestCounter()

	requestID := <-requestIDChan

	pingHandler.sender.SendMsg(
		SenderMessage{
			Proto: &castchannel.CastMessage{
				ProtocolVersion: ProtocolVersion,
				SourceId:        &SenderID,
				DestinationId:   &ReceiverID,
				Namespace:       &NamespaceHeartbeat,
				PayloadType:     PayloadTypeString,
			},
			JsonData: ChromeCastJSONMessage{
				Type:      PayloadTypePing,
				RequestID: requestID,
			},
		})
}

func (pingHandler *PingHandler) GetID() string {
	return pingHandler.id
}

func (pingHandler *PingHandler) OnMsg(msg *castchannel.CastMessage, jsonMsg *ChromeCastJSONMessage) {
	if *msg.Namespace != NamespaceHeartbeat {
		return
	}

	payload := *msg.PayloadUtf8

	if len(payload) == 0 {
		return
	}

	if jsonMsg == nil {
		return
	}

	if jsonMsg.Type != PayloadTypePing {
		return
	}

	requestIDChan := pingHandler.counterGetter.GetRequestCounter()

	requestID := <-requestIDChan

	pingHandler.sender.SendMsg(
		SenderMessage{
			Proto: &castchannel.CastMessage{
				ProtocolVersion: ProtocolVersion,
				SourceId:        &SenderID,
				DestinationId:   &ReceiverID,
				Namespace:       &NamespaceHeartbeat,
				PayloadType:     PayloadTypeString,
			},
			JsonData: ChromeCastJSONMessage{
				Type:      PayloadTypePong,
				RequestID: requestID,
			},
		})
}
