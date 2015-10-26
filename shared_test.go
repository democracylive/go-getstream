package getstream_test

import (
	. "github.com/onedotover/go-getstream"
	"os"
)

const (
	TestToken            = "iFX1l5f_lIUWgZFBnv5UisTTW18"
	TestTargetFeedToken  = "vVm2GeUokcbeFGEPfOWYdbG1ZjY"
	TestTargetFeedToken2 = "EGW6PWbZqmSwYZvxv97-qbPTYas"

	TestFeedSignature        = "flat1 " + TestToken
	TestTargetFeedSignature  = "flat2 " + TestTargetFeedToken
	TestTargetFeedSignature2 = "flat3 " + TestTargetFeedToken2
)

var (
	MockAPIKey    = "api-key"
	MockAPISecret = "gthc2t9gh7pzq52f6cky8w4r4up9dr6rju9w3fjgmkv6cdvvav2ufe5fv7e2r9qy" // From http://getstream.io/docs_rest/#feed_authentication
	MockAppID     = "app-id"

	TestAPIKey    = os.Getenv("GETSTREAM_KEY")
	TestAPISecret = os.Getenv("GETSTREAM_SECRET")
	TestAppID     = os.Getenv("GETSTREAM_APPID")

	TestFeedSlug        = Slug{"flat", "1", ""}
	TestFollowerSlug    = Slug{"flat", "follower", ""}
	TestFollower2Slug   = Slug{"flat", "follower2", ""}
	TestTargetFeedSlug  = Slug{"flat", "target", ""}
	TestTargetFeedSlug2 = Slug{"flat", "target2", ""}
)

// Mock test client for testing signing logic
func MockTestClient(region string) *Client {
	return Connect(MockAPIKey, MockAPISecret, MockAppID, region)
}

// Actual connection to test api integration
func ConnectTestClient(region string) *Client {
	return Connect(TestAPIKey, TestAPISecret, TestAppID, region)
}

func NewSimpleTestActivity() *Activity {
	return &Activity{
		Actor:  TestFeedSlug,
		Verb:   "statusUpdate",
		Object: "user1",
	}
}
func NewTargetTestActivity() *Activity {
	return &Activity{
		Actor:  TestFeedSlug,
		Verb:   "comment",
		Object: "comment1",
		To:     []Slug{TestTargetFeedSlug, TestTargetFeedSlug2},
	}
}
