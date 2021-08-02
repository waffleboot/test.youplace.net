package main

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type cli struct {
	client      *http.Client // вопрос, а может *http.Client
	nextPageUrl string
}

func newCli() (*cli, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c := &cli{client: &http.Client{
		Timeout: 1 * time.Minute,
		Jar:     jar,
	}}
	c.client.CheckRedirect = c.checkRedirect
	return c, nil
}

func (c *cli) checkRedirect(req *http.Request, _ []*http.Request) error {
	c.nextPageUrl = req.URL.String()
	return nil
}

func (c *cli) get(url string) (io.ReadCloser, error) {
	for {
		resp, err := c.client.Get(url) // не очень понятно, в случае 503 вернется ошибка или статус код надо проверять?
		if err != nil {
			return nil, err
		}
		if resp.StatusCode == http.StatusServiceUnavailable {
			continue
		}
		return resp.Body, err
	}
}

func (c *cli) post(data url.Values) (url.Values, error) {
	for {
		resp, err := c.client.PostForm(c.nextPageUrl, data)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusServiceUnavailable {
			continue
		}
		return parseHtml(resp.Body)
	}
}
