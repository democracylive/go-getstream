package getstream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type Client struct {
	http    *http.Client
	baseURL *url.URL // https://api.getstream.io/api/

	key      string
	secret   string
	appID    string
	location string // https://location-api.getstream.io/api/
}

func Connect(key, secret, appID, location string) *Client {
	baseURLStr := "https://api.getstream.io/api/v1.0/"
	if location != "" {
		baseURLStr = "https://" + location + "-api.getstream.io/api/v1.0/"
	}

	baseURL, e := url.Parse(baseURLStr)
	if e != nil {
		panic(e) // failfast, url shouldn't be invalid anyway.
	}

	return &Client{
		http:    &http.Client{},
		baseURL: baseURL,

		key:      key,
		secret:   secret,
		appID:    appID,
		location: location,
	}
}

func (c *Client) BaseURL() *url.URL { return c.baseURL }

func (c *Client) Feed(slug, id string) *Feed {
	return &Feed{
		Client: c,
		slug:   SignSlug(c.secret, Slug{slug, id, ""}),
	}
}

func (c *Client) get(result interface{}, path string, slug Slug, params interface{}) error {
	path, err := addPathParams(path, params)
	if err != nil {
		return err
	}
	return c.request(result, "GET", path, slug, nil)
}

func (c *Client) post(result interface{}, path string, slug Slug, payload interface{}) error {
	return c.request(result, "POST", path, slug, payload)
}

func (c *Client) del(path string, slug Slug) error {
	return c.request(nil, "DELETE", path, slug, nil)
}

func (c *Client) request(result interface{}, method, path string, slug Slug, payload interface{}) error {
	absUrl, e := c.absoluteUrl(path)
	if e != nil {
		return e
	}

	buffer := []byte{}
	if payload != nil {
		if buffer, e = json.Marshal(payload); e != nil {
			return e
		}
	}

	req, e := http.NewRequest(method, absUrl.String(), bytes.NewBuffer(buffer))
	if e != nil {
		return e
	}

	req.Header.Set("Content-Type", "application/json")
	if slug.Token != "" {
		req.Header.Set("Authorization", slug.Signature())
	}

	resp, e := c.http.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()

	buffer, e = ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}

	switch {
	case 200 <= resp.StatusCode && resp.StatusCode < 300: // SUCCESS
		if result != nil {
			if e = json.Unmarshal(buffer, result); e != nil {
				return e
			}
		}

	default:
		err := &Error{}
		if e = json.Unmarshal(buffer, err); e != nil {
			panic(e)
			return errors.New(string(buffer))
		}

		return err
	}

	return nil
}

func (c *Client) absoluteUrl(path string) (result *url.URL, e error) {
	if result, e = url.Parse(path); e != nil {
		return nil, e
	}

	// DEBUG: Use this line to send stuff to a proxy instead.
	// c.baseURL, _ = url.Parse("http://0.0.0.0:8000/")
	result = c.baseURL.ResolveReference(result)

	qs := result.Query()
	qs.Set("api_key", c.key)
	if c.location == "" {
		qs.Set("location", "unspecified")
	} else {
		qs.Set("location", c.location)
	}
	result.RawQuery = qs.Encode()

	return result, nil
}

func addPathParams(path string, params interface{}) (string, error) {
	// Simple conversion from struct to query parameters
	// Only supports basic conversion.
	val := reflect.ValueOf(params)
	if val.IsNil() {
		return path, nil
	}

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return path, errors.New("Params must be a struct.")
	}

	url, err := url.Parse(path)
	if err != nil {
		return path, err
	}

	// Convert the params obejct into query parameters
	query := url.Query()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		typField := typ.Field(i)
		if typField.PkgPath != "" {
			continue // unexported
		}

		tag := typField.Tag.Get("url")
		if tag == "-" {
			continue
		}

		// @todo: Add support for more types or utilize url encoding package like
		// https://github.com/google/go-querystring
		value := ""
		valField := val.Field(i)
		switch valField.Kind() {
		case reflect.String:
			value = valField.String()
		case reflect.Int:
			if valField.Int() != 0 {
				value = fmt.Sprintf("%d", valField)
			}
		default:
			value = fmt.Sprintf("%v", valField)
		}

		// If we have a value set, add it to the query
		if value != "" {
			query.Set(tag, value)
		}
	}
	url.RawQuery = query.Encode()
	return url.String(), nil
}
