package getstream

import (
	"time"
)

type Activity struct {
	Id string `json:"id,omitempty"`
	Actor     string `json:"actor"`
	Verb      string `json:"verb"`
	Object    string `json:"object"`
	Target    string `json:"target,omitempty"`
	RawTime   string `json:"time,omitempty"`
	To        []Slug `json:"to,omitempty"`
	ForeignID string `json:"foreign_id,omitempty"`

	// Response items
	Origin string `json:"origin,omitempty"`
}

func (a *Activity) Sign(secret string) {
	for i := range a.To {
		a.To[i].Sign(secret)
	}
}

type AggregatedActivityBase struct {
	Id string `json:"id"`
	Verb string `json:"verb"`
	Group string `json:"group"`
	ActivityCount int `json:"activity_count"`
	ActorCount int `json:"actor_count"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Activities interface{} `json:"activities"`
}
type AggregatedActivity struct {
	AggregatedActivityBase
	Activities []*Activity `json:"activities"`
}


type ActivitiesResult struct {
	Next        string      `json:"next,omitempty"`
	RawDuration string      `json:"duration,omitempty"`
	Results     interface{} `json:"results,omitempty"`
}

type ActivityOptions struct {
	Limit  int `url:"limit"`
	Offset int `url:"offset"`

	IdGTE string `url:"id_gte"`
	IdGT  string `url:"id_gt"`
	IdLTE string `url:"id_lte"`
	IdLT  string `url:"id_lt"`
//
//	Feeds    []*Feed `json:"feeds"`
//	MarkRead bool    `json:"mark_read"`
//	MarkSeen bool    `json:"mark_seen"`
}

type Notification struct {
	Data    *Update `json"data"`
	Channel string  `json:"channel"`
}

type Update struct {
	Deletes []*Activity
	Inserts []*Activity

	UnreadCount int
	UnseenCount int
	PublishedAt time.Time
}

type FollowPost struct {
	Target            string `json:"target"`
	ActivityCopyLimit int    `json:"activity_copy_limit,omitempty"`
}

type FollowPostResult struct {
	Duration string `json:"duration"`
}

type FollowingOptions struct {
	Limit  int    `url:"limit"`
	Offset int    `url:"offset"`
	Filter string `url:"filter"`
}

type FollowingResult struct {
	RawDuration string           `json:"duration,omitempty"`
	Results     []*FollowingInfo `json:"results,omitempty"`
}

type FollowingInfo struct {
	FeedId    string    `json:"feed_id"`
	TargetId  string    `json:"target_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
