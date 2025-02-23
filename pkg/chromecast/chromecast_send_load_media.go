package chromecast

import (
	"log"

	"github.com/vjerci/gochromecast/pkg/chromecast/proto/castchannel"
)

var PayloadTypeLoad = "LOAD"

type PayloadLoad struct {
	Type           string            `json:"type"`
	Media          *PayloadLoadMedia `json:"media"`
	CurrentTime    int               `json:"currentTime"`
	ActiveTrackIDs []int             `json:"activeTrackIds"`
	RequestID      int               `json:"requestId"`
}

type PayloadLoadMedia struct {
	ContentID   string               `json:"contentId"`
	StreamType  string               `json:"streamType"` // can be ignored or LIVE
	ContentType string               `json:"contentType"`
	Tracks      []*PayloadLoadTracks `json:"tracks"`
}

type PayloadLoadTracks struct {
	TrackID          int    `json:"trackId"`
	TrackContentID   string `json:"trackContentId"`
	TrackContentType string `json:"trackContentType"`
	Type             string `json:"type"`
	SubType          string `json:"subtype"`
	Language         string `json:"language"`
	Name             string `json:"name"`
}

var NamespaceMedia = "urn:x-cast:com.google.cast.media"

func (client *Client) SendLoadMedia(requestCounter *RequestCounter, req PlayMediaRequest, transportID string, sender *Sender) {
	requestIDCounter := <-requestCounter.GetRequestCounter()

	loadJSONmsg := PayloadLoad{
		Type: PayloadTypeLoad,
		// CurrentTime: 0, // 3* 60 to play 3rd minute
		Media: &PayloadLoadMedia{
			ContentID:   req.MediaURL,
			ContentType: "application/x-mpegurl",
		},
		RequestID: requestIDCounter,
	}

	log.Printf("chromecast subtitle url is %s", req.SubtitlesURL)

	if req.SubtitlesURL != "" {
		loadJSONmsg.ActiveTrackIDs = []int{3}
		loadJSONmsg.Media.Tracks = []*PayloadLoadTracks{
			{
				TrackID:          3,
				TrackContentID:   req.SubtitlesURL,
				TrackContentType: "text/vtt",
				Type:             "TEXT",
				SubType:          "SUBTITLES",
				Language:         "en",
				Name:             "en subtitles",
			},
		}
	}

	sender.SendMsg(SenderMessage{
		Proto: &castchannel.CastMessage{
			ProtocolVersion: ProtocolVersion,
			SourceId:        &SenderID,
			DestinationId:   &transportID,
			Namespace:       &NamespaceMedia,
			PayloadType:     PayloadTypeString,
		},
		JsonData: loadJSONmsg,
	})
}
