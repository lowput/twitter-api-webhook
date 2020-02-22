package twitter_api_webhook

import (
	"cloud.google.com/go/storage"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type HmacResult struct {
	ResponseToken string	`json:"response_token"`
}

type TweetCreateEvent struct {
	ForUseId string 		`json:"for_user_id"`
	Event	[]*CreateEvent	`json:"tweet_create_events"`
}

type CreateEvent struct {
	Entities Entities `json:"entities"`
}

type Entities struct {
	MediaEntities []*MediaEntities	`json:"media"`
}

type MediaEntities struct {
	Indices     [2]int `json:"indices"`
	DisplayURL  string  `json:"display_url"`
	ExpandedURL string  `json:"expanded_url"`
	URL         string  `json:"url"`
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
			save(media.MediaURLHttps + ":large", media.IDStr)
		}
	}
}

func save(url string, id string)  {
	if len(url) == 0 {
		log.Fatal("empty url")
		return
	} else {
		log.Printf("Backet save : %s", url)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	name := fmt.Sprintf("%s.jpg", id)
	bucket := client.Bucket(os.Getenv("BUCKET_NAME")).Object(name).NewWriter(ctx)
	defer bucket.Close()

	res , err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	io.Copy(bucket, res.Body)
}
