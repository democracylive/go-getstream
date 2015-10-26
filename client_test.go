package getstream_test

import (
	a "github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_BaseURL(t *testing.T) {
	locations := map[string]string{
		"":        "https://api.getstream.io/api/v1.0/",
		"us-east": "https://us-east-api.getstream.io/api/v1.0/",
		"xyz":     "https://xyz-api.getstream.io/api/v1.0/",
	}

	for location, url := range locations {
		client := MockTestClient(location)
		a.Equal(t, url, client.BaseURL().String())
	}
}

func TestClient_Feed(t *testing.T) {
	client := MockTestClient("")
	feed := client.Feed(TestFeedSlug.Slug, TestFeedSlug.ID)

	slug := TestFeedSlug
	slug.Sign(MockAPISecret)

	a.Equal(t, slug, feed.Slug())
	a.Equal(t, TestFeedSignature, feed.Slug().Signature())
}
