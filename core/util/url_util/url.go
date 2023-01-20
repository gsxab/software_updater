package url_util

import "net/url"

func RelativeURL(baseURL, relativeURL string) (*url.URL, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	rel, err := base.Parse(relativeURL)
	if err != nil {
		return nil, err
	}
	return rel, nil
}
