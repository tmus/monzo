package monzo

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type FeedItemType string

const FeedItemBasic FeedItemType = "basic"

type FeedItem struct {
	itemType FeedItemType
	title    string
	body     string
	imageURL string

	bgColor    string
	titleColor string
	bodyColor  string
}

func MakeBasicFeedItem(title string, body string) FeedItem {
	return FeedItem{
		itemType: FeedItemBasic,
		title:    title,
		body:     body,
		// If an imageURL is not set, Monzo complains about a non-URL
		// value, so just set it to something random here. If the
		// image is not found, it's just not displayed at all.
		imageURL: "https://tomm.us/",
	}
}

func (a Account) AddFeedItem(fi FeedItem) error {
	endpoint := "/feed"
	data := url.Values{}
	data.Add("account_id", a.ID)
	data.Add("type", string(fi.itemType))
	data.Add("params[title]", fi.title)
	data.Add("params[body]", fi.body)
	data.Add("params[image_url]", fi.imageURL)

	req, err := a.client.NewRequest(http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := a.client.Do(req)

	b := new(bytes.Buffer)
	b.ReadFrom(resp.Body)
	str := b.String()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to add metadata: %s", str)
	}

	return nil
}
