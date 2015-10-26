package getstream

type Feed struct {
	*Client
	slug Slug
}

func (f *Feed) Slug() Slug { return f.slug }

func (f *Feed) AddActivity(activity *Activity) (*Activity, error) {
	activity = SignActivity(f.secret, activity)

	result := &Activity{}
	e := f.post(result, f.url(), f.slug, activity)
	return result, e
}

func (f *Feed) AddActivities(activities []*Activity) error {
	signeds := make([]*Activity, len(activities), len(activities))
	for i, activity := range activities {
		signeds[i] = SignActivity(f.secret, activity)
	}

	// TODO: A result type to recieve the listing result.
	panic("not yet implemented.")
}

func (f *Feed) Activities(opt *Options) ([]*Activity, error) {
	result := ActivitiesResult{}
	e := f.get(&result, f.url(), f.slug, opt)
	return result.Results, e
}

func (f *Feed) RemoveActivity(id string) error {
	return f.del(f.url()+id+"/", f.slug)
}

func (f *Feed) Follow(feed, id string) error {
	result := FollowPostResult{}
	target := feed + ":" + id
	return f.post(&result, f.url()+"following/", f.slug, FollowPost{Target: target})
}

func (f *Feed) Unfollow(feed, id string) error {
	target := feed + ":" + id
	return f.del(f.url()+"following/"+target+"/", f.slug)
}

func (f *Feed) Followers(opt *FollowingOptions) ([]*FollowingInfo, error) {
	result := FollowingResult{}
	err := f.get(&result, f.url()+"followers/", f.slug, opt)
	return result.Results, err
}

func (f *Feed) Following(opt *FollowingOptions) ([]*FollowingInfo, error) {
	result := FollowingResult{}
	err := f.get(&result, f.url()+"following/", f.slug, opt)
	return result.Results, err
}

func (f *Feed) url() string {
	return "feed/" + f.slug.Slug + "/" + f.slug.ID + "/"
}
