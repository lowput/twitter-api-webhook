package twitter_api_webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type HmacResult struct {
	ResponseToken string `json:"response_token"`
}

type TweetCreateEvent struct {
	ForUseId string         `json:"for_user_id"`
	Event    []*CreateEvent `json:"tweet_create_events"`
}

type CreateEvent struct {
	Entities Entities `json:"entities"`
}

type Entities struct {
	MediaEntities []*MediaEntities `json:"media"`
}

type MediaEntities struct {
	Indices           [2]int     `json:"indices"`
	DisplayURL        string     `json:"display_url"`
	ExpandedURL       string     `json:"expanded_url"`
	URL               string     `json:"url"`
	ID                int64      `json:"id"`
	IDStr             string     `json:"id_str"`
	MediaURL          string     `json:"media_url"`
	MediaURLHttps     string     `json:"media_url_https"`
	SourceStatusID    int64      `json:"source_status_id"`
	SourceStatusIDStr string     `json:"source_status_id_str"`
	Type              string     `json:"type"`
	Sizes             MediaSizes `json:"sizes"`
}

type MediaSize struct {
	Width  int    `json:"w"`
	Height int    `json:"h"`
	Resize string `json:"resize"`
}

type MediaSizes struct {
	Thumb  MediaSize `json:"thumb"`
	Large  MediaSize `json:"large"`
	Medium MediaSize `json:"medium"`
	Small  MediaSize `json:"small"`
}

func TwitterMediaImageSave(body []byte) {
	var tweet TweetCreateEvent
	if err := json.Unmarshal(body, &tweet); err != nil {
		log.Fatal(err)
	}
	if len(tweet.Event) == 0 {
		return
	}

	for _, event := range tweet.Event {
		for _, media := range event.Entities.MediaEntities {
			slackPost(media.MediaURLHttps+":large", fmt.Sprintf("%s.jpg", media.IDStr))
		}
	}
}

func slackPost(imageUrl string, fileName string) {
	if len(imageUrl) == 0 {
		log.Fatal("empty url")
		return
	}

	log.Printf("Post slack : %s", imageUrl)

	image, err := http.Get(imageUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer image.Body.Close()

	// Slack Post
	url := "https://slack.com/api/files.upload?token=" + os.Getenv("SLACK_ACCESS_TOKEN") + "&channels=" + os.Getenv("SLACK_CHANNEL_ID")
	body := &bytes.Buffer{}

	mw := multipart.NewWriter(body)
	fw, _ := mw.CreateFormFile("file", fileName)
	io.Copy(fw, image.Body)

	contentType := mw.FormDataContentType()
	mw.Close()

	res, err := http.Post(url, contentType, body)
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()
}
