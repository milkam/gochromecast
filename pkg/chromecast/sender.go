package chromecast

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"log"

	"github.com/milkam/gochromecast/pkg/chromecast/proto/castchannel"
	"google.golang.org/protobuf/proto"
)

type Sender struct {
	ctx  context.Context
	conn *tls.Conn

	queue chan SenderMessage
}

type SenderMessage struct {
	Proto    *castchannel.CastMessage
	JsonData any
}

func NewSender(ctx context.Context, conn *tls.Conn) *Sender {
	return &Sender{
		ctx:   ctx,
		conn:  conn,
		queue: make(chan SenderMessage, 100),
	}
}

func (sender *Sender) Start() {
	go sender.listen()
}

func (sender *Sender) SendMsg(msg SenderMessage) {
	go func() {
		sender.queue <- msg
	}()
}

func (sender *Sender) listen() {
	for {
		select {
		case <-sender.ctx.Done():
			return
		case msg := <-sender.queue:
			sender.sendMsgOverTLS(msg)
		}
	}
}

func (sender *Sender) sendMsgOverTLS(msg SenderMessage) {
	// if there is msg payload marshal it into proto msg
	if msg.JsonData != nil {
		jsonMsgBytes, err := json.Marshal(&msg.JsonData)
		if err != nil {
			log.Printf("failed to marshal json message for sending '%s' %#v", err, msg.JsonData)
			return
		}

		jsonMsgBytesS := string(jsonMsgBytes)

		msg.Proto.PayloadUtf8 = &jsonMsgBytesS
	}

	log.Printf("sending message %s", msg.Proto.String())

	castMsgBytes, err := proto.Marshal(msg.Proto)
	if err != nil {
		log.Printf("failed to marshal proto msg for sending '%s' '%#v'", err, msg.Proto)
		return
	}

	// specify it is big endian
	protoLength := make([]byte, 4)

	binary.BigEndian.PutUint32(protoLength, uint32(len(castMsgBytes)))

	fullData := append([]byte{}, protoLength...)
	fullData = append(fullData, castMsgBytes...)

	_, err = sender.conn.Write(fullData)
	if err != nil {
		log.Printf("failed to send marshaled message over conn '%s'", err)
	}
}
