package chromecast

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/vjerci/gochromecast/pkg/chromecast/proto/castchannel"
)

type MessageSubscriberUnsubscribe interface {
	Unsubscribe(id string)
}

type Message struct {
	Json  *ChromeCastJSONMessage
	Proto *castchannel.CastMessage
}

type MessageSubscriber struct {
	id  string
	ctx context.Context

	C chan *Message

	config *MessageSubscriberConfig
}

type MessageSubscriberConfig struct {
	UnSub MessageSubscriberUnsubscribe

	TargetNamespace   string
	TargetPayloadType string
	TargetReceiver    string
	TargetSender      string

	TargetRequestID int

	Debug bool
}

func NewMessageSubscriber(ctx context.Context, config *MessageSubscriberConfig) *MessageSubscriber {
	return &MessageSubscriber{
		id:  uuid.NewString(),
		ctx: ctx,
		C:   make(chan *Message),

		config: config,
	}
}

func (msgSub *MessageSubscriber) GetID() string {
	return msgSub.id
}

func (msgSub *MessageSubscriber) OnMsg(protoMsg *castchannel.CastMessage, jsonMessage *ChromeCastJSONMessage) {
	log.Printf("got message %s", protoMsg)

	if msgSub.config.Debug {
		log.Printf("got msg with payload %s", *protoMsg.PayloadUtf8)
	}

	if msgSub.config.TargetNamespace != *protoMsg.Namespace {
		if msgSub.config.Debug {
			log.Printf("namespace doesn't match '%s' %s'", msgSub.config.TargetNamespace, *protoMsg.Namespace)
		}

		return
	}

	if jsonMessage == nil {
		if msgSub.config.Debug {
			log.Printf("jsonMessage is nil %#v", jsonMessage)
		}
		return
	}

	if msgSub.config.TargetPayloadType != jsonMessage.Type {
		if msgSub.config.Debug {
			log.Printf("payload type doesn't match '%s' '%s'", msgSub.config.TargetPayloadType, jsonMessage.Type)
		}
		return
	}

	if msgSub.config.TargetSender != *protoMsg.SourceId {
		if msgSub.config.Debug {
			log.Printf("source id  doesn't match '%s' '%s'", msgSub.config.TargetSender, *protoMsg.SourceId)
		}
		return
	}

	if msgSub.config.TargetReceiver != *protoMsg.DestinationId && *protoMsg.DestinationId != "*" {
		if msgSub.config.Debug {
			log.Printf("destionation id  doesn't match '%s' '%s'", msgSub.config.TargetReceiver, *protoMsg.DestinationId)
		}
		return
	}

	if msgSub.config.TargetRequestID != jsonMessage.RequestID {
		if msgSub.config.Debug {
			log.Printf("request id  doesn't match '%s' '%s'", msgSub.config.TargetRequestID, jsonMessage.RequestID)
		}
		return
	}

	if msgSub.config.Debug {
		log.Printf("resolving subscription")
	}

	msgSub.C <- &Message{
		Proto: protoMsg,
		Json:  jsonMessage,
	}

	msgSub.config.UnSub.Unsubscribe(msgSub.id)
}
