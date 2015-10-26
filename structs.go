package getstream

import (
	"time"
)

type Activity struct {
	ID        string `json:"id,omitempty"`
	Actor     Slug   `json:"actor"`
	Verb      string `json:"verb"`
	Object    string `json:"object"`
	Target    string `json:"target,omitempty"`
	RawTime   string `json:"time,omitempty"`
	To        []Slug `json:"to,omitempty"`
	ForeignID string `json:"foreign_id,omitempty"`
}

type ActivitiesResult struct {
	Next        string      `json:"next,omitempty"`
	RawDuration string      `json:"duration,omitempty"`
	Results     []*Activity `json:"results,omitempty"`
}

type Options struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`

	IdGTE string `json:"id_gte"`
	IdGT  string `json:"id_gt"`
	IdLTE string `json:"id_lte"`
	IdLT  string `json:"id_lt"`

	Feeds    []*Feed `json:"feeds"`
	MarkRead bool    `json:"mark_read"`
	MarkSeen bool    `json:"mark_seen"`
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
