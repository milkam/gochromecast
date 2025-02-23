package chromecast

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

var NamespaceReceiver = "urn:x-cast:com.google.cast.receiver"

const PayloadTypeRecieverStatus = "RECEIVER_STATUS"

func (client *Client) SubscribeRecieverStatus(ctx context.Context, tlsReader *Reader, requestID int) *MessageSubscriber {
	statusSub := NewMessageSubscriber(ctx, &MessageSubscriberConfig{
		UnSub:             tlsReader,
		TargetNamespace:   NamespaceReceiver,
		TargetPayloadType: PayloadTypeRecieverStatus,
		TargetReceiver:    SenderID,
		TargetSender:      ReceiverID,
		TargetRequestID:   requestID,
		Debug:             true,
	})
	tlsReader.Subscribe(statusSub)

	return statusSub
}

var ErrPlayMediaTimeoutStatus = errors.New("getting media status timed out")

func (client *Client) WaitForTransportID(statusSub *MessageSubscriber) (transportID string, err error) {
	timeout := time.NewTimer(15 * time.Second)

	// wait for status to show which contains our transport id
	var msg *Message

	select {
	case <-timeout.C:
		return "", errors.Join(ErrPlayMediaTimeoutStatus, err)
	case msg = <-statusSub.C:
	}

	transportID, err = getTransportID([]byte(*msg.Proto.PayloadUtf8))
	if err != nil {
		return "", errors.Join(ErrPlayMediaTransportID, err)
	}

	return transportID, nil
}

var ErrGetTransportIDJson = errors.New("json unmarshal failed")
var ErrGetTransportIDNoApp = errors.New("app with default id not found")

type ReceiverStatus struct {
	Status struct {
		Applications []struct {
			AppId          string `json:"appId"`
			TransportId    string `json:"transportId"`
			UniversalAppId string `json:"universalAppId"`
		} `json:"applications"`
	} `json:"status"`
}

func getTransportID(data []byte) (string, error) {
	var receiverStatus ReceiverStatus
	err := json.Unmarshal(data, &receiverStatus)
	if err != nil {
		return "", errors.Join(ErrGetTransportIDJson, err)
	}

	for _, app := range receiverStatus.Status.Applications {
		if app.AppId == DefaultMediaAppID {
			return app.TransportId, nil
		}
	}

	return "", ErrGetTransportIDNoApp
}
