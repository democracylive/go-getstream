package getstream_test

import (
	"github.com/onedotover/go-getstream"
	a "github.com/stretchr/testify/assert"
	"testing"
	"log"
)

func TestFeed(t *testing.T) {
	client := ConnectTestClient("")
	feed := client.Feed(TestFeedSlug.Slug, TestFeedSlug.ID)

	a.Equal(t, TestFeedSlug.Slug + ":" + TestFeedSlug.ID, feed.Id())
}
func TestActivityCrud(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := ConnectTestClient("")
	feed := client.Feed(TestFeedSlug.Slug, TestFeedSlug.ID)
	activity := TestActivityTarget

	t.Log("adding activity...")
	e := feed.AddActivity(&activity)
	a.NoError(t, e)
	a.NotEmpty(t, activity.Id)

	t.Log("listing added activities...")
	activities := []*getstream.Activity{}
	_, e = feed.GetActivities(&activities, nil)
	a.NoError(t, e)
	a.NotEmpty(t, activities)
	a.Len(t, activities, 1) // otherwise we might be getting result from another test run.
	a.Equal(t, activity.Id, activities[0].Id)

//	t.Log("removing added activity...")
//	e = feed.RemoveActivity(activity.Id)
//	a.NoError(t, e)
//
//	t.Log("listing added activities again...")
//	activities = []*getstream.Activity{}
//	_, e = feed.GetActivities(&activities, nil)
//	a.NoError(t, e)
//	a.Empty(t, activities)
}
func TestForeignId(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := ConnectTestClient("")
	feed := client.Feed(TestFeedSlug.Slug, TestFeedSlug.ID)

	// Create a simple test
	activity := TestActivitySimple
	feed.AddActivity(&activity)

	// Delete it
	feed.RemoveActivityByForeignId(activity.ForeignID)

	// See if it deleted
	activities := []getstream.Activity{}
	feed.GetActivities(&activities, nil)
	a.Empty(t, activities)
}
func TestFollow(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := ConnectTestClient("")

	// Check the follower feed and make sure it's empty
	follower := client.Feed(TestFollowerSlug.Slug, TestFollowerSlug.ID)
	activities := []*getstream.Activity{}
	_, err := follower.GetActivities(&activities, nil)
	a.NoError(t, err)
	if !a.Len(t, activities, 0) {
		return // The rest of this test depends on just one activity
	}

	// Post an activity to the main test feed
	feed := client.Feed(TestFeedSlug.Slug, TestFeedSlug.ID)
	activity := TestActivitySimple
	err = feed.AddActivity(&activity)
	a.NoError(t, err)

	// Now follow the main feed
	err = follower.Follow(TestFeedSlug.Slug, TestFeedSlug.ID)
	a.NoError(t, err)

	// We should now have content in the follower feed
	activities = []*getstream.Activity{}
	_, err = follower.GetActivities(&activities, nil)
	a.NoError(t, err)
	if a.Len(t, activities, 1) {
		a.Equal(t, activity.Id, activities[0].Id)
	}

	// Add a second follower so we can test limits and paging
	follower2 := client.Feed(TestFollower2Slug.Slug, TestFollower2Slug.ID)
	follower2.Follow(TestFeedSlug.Slug, TestFeedSlug.ID)

	// The main test feed should have two followers now
	followers, err := feed.Followers(nil)
	a.NoError(t, err)
	if a.Len(t, followers, 2) {
		a.Equal(t, follower2.Id(), followers[0].FeedId)
		a.Equal(t, follower.Id(), followers[1].FeedId)
	}

	// And the follower should be following one feed (the main test feed)
	following, err := follower.Following(nil)
	a.NoError(t, err)
	if a.Len(t, following, 1) {
		a.Equal(t, feed.Id(), following[0].TargetId)
	}

	// Test following options
	opt := &getstream.FollowingOptions{Limit: 1}
	followers, err = feed.Followers(opt)
	a.NoError(t, err)
	if a.Len(t, followers, 1) {
		a.Equal(t, follower2.Id(), followers[0].FeedId)
	}

	opt = &getstream.FollowingOptions{Limit: 1, Offset: 1}
	followers, err = feed.Followers(opt)
	a.NoError(t, err)
	if a.Len(t, followers, 1) {
		a.Equal(t, follower.Id(), followers[0].FeedId)
	}

	opt = &getstream.FollowingOptions{Filter: feed.Id()}
	following, err = follower2.Following(opt)
	a.NoError(t, err)
	if a.Len(t, following, 1) {
		a.Equal(t, feed.Id(), following[0].TargetId)
	}

	opt = &getstream.FollowingOptions{Filter: "nonexistent:id"}
	following, err = follower2.Following(opt)
	a.NoError(t, err)
	a.Len(t, following, 0)

	// Test convenience functions
	a.True(t, follower.IsFollowing(TestFeedSlug.Slug, TestFeedSlug.ID))
	a.False(t, follower2.IsFollowing(TestFeedSlug.Slug, TestFeedSlug.ID))

	// Unfollow
	err = follower.Unfollow(TestFeedSlug.Slug, TestFeedSlug.ID)
	a.NoError(t, err)

	err = follower2.Unfollow(TestFeedSlug.Slug, TestFeedSlug.ID)
	a.NoError(t, err)

	// We should be left with a empty feeds
	activities = []*getstream.Activity{}
	_, err = follower.GetActivities(&activities, nil)
	a.NoError(t, err)
	a.Len(t, activities, 0)

	// And we should be left with no following ties
	followers, err = feed.Followers(nil)
	a.NoError(t, err)
	a.Len(t, followers, 0)

	following, err = follower.Following(nil)
	a.NoError(t, err)
	a.Len(t, following, 0)

	following, err = follower2.Following(nil)
	a.NoError(t, err)
	a.Len(t, following, 0)

	// Clean up the post in the main test feed
	err = feed.RemoveActivity(activity.Id)
	a.NoError(t, err)
}

func TestExtendedActivity(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := ConnectTestClient("")
	feed := client.Feed(TestFeedSlug.Slug, TestFeedSlug.ID)

	activity := TestActivityExtended
	err := feed.AddActivity(&activity)
	a.NoError(t, err)

	activities := []ExtendedActivity{}
	_, err = feed.GetActivities(&activities, nil)

	a.NoError(t, err)
	if a.Len(t, activities, 1) {
		a.Equal(t, activity.Id, activities[0].Id)
		a.Equal(t, activity.Title, activities[0].Title)
	}

	feed.RemoveActivity(activity.Id)
}

func TestAggregated(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	client := ConnectTestClient("")
	feed := client.Feed(TestFeedSlug.Slug, TestFeedSlug.ID)
	aggregated := client.Feed(TestAggregatedSlug.Slug, TestAggregatedSlug.ID)
	aggregated.Follow(TestFeedSlug.Slug, TestFeedSlug.ID)

	activity := TestActivityExtended
	feed.AddActivity(&activity)
	log.Println("activity", activity, TestActivityExtended)
	activity2 := TestActivityExtended
	activity2.Comment = "Commenting again"
	feed.AddActivity(&activity2)
	log.Println("added", activity, activity2)

	results := []AggregatedExtendedActivity{}
	_, err := aggregated.GetActivities(&results, nil)
	a.NoError(t, err)
	a.Equal(t, results[0].ActorCount, 1)
	a.Equal(t, results[0].Verb, activity.Verb)
	if a.Len(t, results, 1) {// Should have one aggreggated activity
		if a.Len(t, results[0].Activities, 2) { // With two activities in it
			a.Equal(t, activity2.Comment, results[0].Activities[0].Comment)
			a.Equal(t, activity.Comment, results[0].Activities[1].Comment)
		}
	}

	aggregated.Unfollow(TestFeedSlug.Slug, TestFeedSlug.ID)
	feed.RemoveActivity(activity.Id)
	feed.RemoveActivity(activity2.Id)
}