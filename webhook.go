package twitter_api_webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func TwitterApi(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		uri, _ := url.Parse(r.RequestURI)
		query := uri.Query()
		crc := query["crc_token"][0]
		hmacRes := HmacResult{"sha256=" + makeHmac(os.Getenv("TWITTER_CONSUMER_SECRET"), crc)}

		if res, err := json.Marshal(hmacRes); err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(res)
		}
	case "POST":
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			TwitterMediaImageSave(body)
		}
	}
}

// Make HMAC and Base64 encode.
func makeHmac(key, msg string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	hashBinary := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashBinary)
}
