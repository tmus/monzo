package monzo

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// FeedItemType defines the type of FeedItem to create.
type FeedItemType string

// FeedItemBasic is currently the only option supported by
// Monzo for use through the API.
const FeedItemBasic FeedItemType = "basic"

// FeedItem represents a single row on the Monzo feed. API-
// generated feed items are dismissable by the user and
// should not be considered permanent.
type FeedItem struct {
	itemType FeedItemType
	title    string
	body     string
	imageURL string

	bgColor    string
	titleColor string
	bodyColor  string
}

// MakeFeedItem creates a basic item to display in the Monzo feed.
// Additional values can be chained to the feed item to customise
// it further.
func MakeFeedItem(title string, body string) *FeedItem {
	return &FeedItem{
		itemType: FeedItemBasic,
		title:    title,
		body:     body,
		// If an imageURL is not set, Monzo complains about a non-URL
		// value, so just set it to something random here. If the
		// image is not found, it's just not displayed at all.
		imageURL: "https://tomm.us/",
	}
}

// Image sets the image url on a FeedItem.
func (fi *FeedItem) Image(url string) *FeedItem {
	fi.imageURL = url
	return fi
}

// BackgroundColor sets the color of the background.
func (fi *FeedItem) BackgroundColor(hex string) *FeedItem {
	fi.bgColor = hex
	return fi
}

// TitleColor sets the color of the title text.
func (fi *FeedItem) TitleColor(hex string) *FeedItem {
	fi.titleColor = hex
	return fi
}

// BodyColor changes the color of the body text.
func (fi *FeedItem) BodyColor(hex string) *FeedItem {
	fi.bodyColor = hex
	return fi
}

// AddFeedItem adds a passed FeedItem to the Monzo Account to
// display in the feed.
func (a Account) AddFeedItem(fi *FeedItem) error {
	endpoint := "/feed"
	data := url.Values{}
	data.Add("account_id", a.ID)
	data.Add("type", string(fi.itemType))
	data.Add("params[title]", fi.title)
	data.Add("params[body]", fi.body)
	data.Add("params[image_url]", fi.imageURL)

	if fi.bgColor != "" {
		data.Add("params[background_color]", fi.bgColor)
	}

	if fi.titleColor != "" {
		data.Add("params[title_color]", fi.titleColor)
	}

	if fi.bodyColor != "" {
		data.Add("params[body_color]", fi.bodyColor)
	}

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
