package getstream_test

import (
	. "github.com/onedotover/go-getstream"
	"os"
)

const (
	TestToken            = "iFX1l5f_lIUWgZFBnv5UisTTW18"
	TestTargetFeedToken  = "Gn3Dzit2bk3MuClNjvbXO-49RMI"
	TestTargetFeedToken2 = "oRJR8IFZU-QiBX89CSuPcwqxz8I"

	TestFeedSignature        = "flat1 " + TestToken
	TestTargetFeedSignature  = "flattarget " + TestTargetFeedToken
	TestTargetFeedSignature2 = "flattarget2 " + TestTargetFeedToken2
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

	TestAggregatedSlug = Slug{"aggregated", "1", ""}

	TestActivitySimple = Activity{
		Actor:  "user1",
		Verb:   "statusUpdate",
		Object: "user1",
		ForeignID: "user1",
	}

	TestActivityTarget = Activity{
		Actor:  "user1",
		Verb:   "taret",
		Object: "user1",
		To:     []Slug{TestTargetFeedSlug, TestTargetFeedSlug2},
	}

	TestActivityExtended = ExtendedActivity{
		Activity: Activity{
			Actor:  "user3",
			Verb:   "comment",
			Object: "article1",
		},
		Title:   "Comment title",
		Comment: "Comment on this item.",
	}
)

// Mock test client for testing signing logic
func MockTestClient(region string) *Client {
	return Connect(MockAPIKey, MockAPISecret, MockAppID, region)
}

// Actual connection to test api integration
func ConnectTestClient(region string) *Client {
	return Connect(TestAPIKey, TestAPISecret, TestAppID, region)
}

type ExtendedActivity struct {
	Activity
	Title   string `json:"title"`
	Comment string `json:"comment"`
}

type AggregatedExtendedActivity struct {
	AggregatedActivityBase
	Activities []ExtendedActivity `json:"activities"`
}