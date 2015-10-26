package getstream
import (
	"reflect"
	"errors"
)

type Feed struct {
	*Client
	slug Slug
}

func (f *Feed) Slug() Slug { return f.slug }

// Wrapper can either be a *getstream.Activity{} or
// it can be a struct that extends getstream.Activity like this:
// type ExtendedActivity struct {
//   *getstream.Activity
//   AdditionalField string `json:"additionalField"`
// }
func (f *Feed) AddActivity(wrapper interface{}) error {
	val := reflect.ValueOf(wrapper)
	if val.Kind() != reflect.Ptr {
		return errors.New("Must pass a pointer to an activity into AddActivity()")
	}

	// Allow for an Activity to be passed in or a struct that extends Activity
//	var activity *Activity
	if val.Elem().Type().String() == "getstream.Activity" {
		wrapper.(*Activity).Sign(f.secret)
	} else {
		elem := val.Elem()
		field := elem.FieldByName("Activity")
		if !field.IsValid() {
			return errors.New("Activity must extend getstream.Activity")
		}
		if field.Kind() != reflect.Ptr {
			return errors.New("Anonyous field *Activity must be a pointer")
		}
		activity := field.Interface().(*Activity)
		activity.Sign(f.secret)
//		field.Elem().Set(reflect.ValueOf(activity).Elem())
	}

//	result := reflect.New(reflect.TypeOf(wrapper).Elem()).Interface()
	e := f.post(wrapper, f.url(), f.slug, wrapper)
	return e
}

func (f *Feed) AddActivities(activities []*Activity) error {
	for i := range activities {
		activities[i].Sign(f.secret)
	}

	// TODO: A result type to recieve the listing result.
	panic("not yet implemented.")
}

func (f *Feed) Activities(target interface{}, opt *Options) (string, error) {
	result := ActivitiesResult{}
	result.Results = target
	e := f.get(&result, f.url(), f.slug, opt)
	return result.Next, e
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
