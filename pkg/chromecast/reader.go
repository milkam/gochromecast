package chromecast

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/vjerci/gochromecast/pkg/chromecast/proto/castchannel"
	"google.golang.org/protobuf/proto"
)

type ReaderSubscriber interface {
	GetID() string
	OnMsg(*castchannel.CastMessage, *ChromeCastJSONMessage)
}

type ChromeCastJSONMessage struct {
	Type      string `json:"type"`
	RequestID int    `json:"requestId"`
}

type Reader struct {
	ctx              context.Context
	closed           bool
	conn             *tls.Conn
	subscribers      []ReaderSubscriber
	addSubscriber    chan ReaderSubscriber
	removeSubscriber chan string
}

func NewReader(ctx context.Context, conn *tls.Conn) *Reader {
	r := &Reader{
		ctx:              ctx,
		closed:           false,
		conn:             conn,
		subscribers:      make([]ReaderSubscriber, 0),
		addSubscriber:    make(chan ReaderSubscriber, 100),
		removeSubscriber: make(chan string, 100),
	}

	return r
}

func (reader *Reader) Start() {
	go reader.read()
	go reader.subscribeHandler()
}

func (reader *Reader) Close() {
	reader.closed = true
}

func (reader *Reader) Subscribe(sub ReaderSubscriber) {
	go func() {
		reader.addSubscriber <- sub
	}()
}

func (reader *Reader) Unsubscribe(id string) {
	go func() {
		reader.removeSubscriber <- id
	}()
}

func (reader *Reader) subscribeHandler() {
	for {
		select {
		case <-reader.ctx.Done():
			reader.subscribers = nil
			return
		case sub := <-reader.addSubscriber:
			reader.subscribers = append(reader.subscribers, sub)
		case unsub := <-reader.removeSubscriber:
			cleanedSubs := []ReaderSubscriber{}
			for _, sub := range reader.subscribers {
				if sub.GetID() != unsub {
					cleanedSubs = append(cleanedSubs, sub)
				}
			}
			reader.subscribers = cleanedSubs
		}
	}
}

var ErrReaderRead = errors.New("reading tls failed")
var ErrReaderUnmarshal = errors.New("unmarshal proto failed")

func (reader *Reader) read() {
	for {
		if reader.closed {
			log.Printf("tls reader closed closing loop")
			return
		}

		data := make([]byte, 10000)
		i, err := reader.conn.Read(data)
		if err != nil {
			log.Printf("reading conn failed: %s", errors.Join(ErrReaderRead, err))
			return
		}

		if i == 0 {
			log.Printf("read 0 bytes sleeping 200 ms")
			time.Sleep(200 * time.Millisecond)
			continue
		}

		data = data[4:i]
		var portoMsg castchannel.CastMessage
		err = proto.Unmarshal(data, &portoMsg)
		if err != nil {
			log.Printf("unmarshaling proto msg failed %s", errors.Join(ErrReaderUnmarshal, err))
			return
		}

		log.Printf("---msRecieved")
		log.Printf("%s", portoMsg.String())

		jsonMessageS := *portoMsg.PayloadUtf8

		if len(jsonMessageS) != 0 {
			var jsonMsg ChromeCastJSONMessage

			err := json.Unmarshal([]byte(jsonMessageS), &jsonMsg)
			if err != nil {
				log.Print("failed to unmarshal json payload in ping handler")
				continue
			}

			for _, subscriber := range reader.subscribers {
				go subscriber.OnMsg(&portoMsg, &jsonMsg)
			}

			continue
		}

		for _, subscriber := range reader.subscribers {
			go subscriber.OnMsg(&portoMsg, nil)
		}
	}
}
