package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var rejected = errors.New("rejected")

const mainUrl = "http://test.youplace.net"

type cli struct {
	client *http.Client // вопрос, а может *http.Client
	url    string
}

func (c *cli) saveRedirect(req *http.Request, _ []*http.Request) error {
	c.url = req.URL.String()
	return nil
}

func newCli() (*cli, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	c := &cli{client: &http.Client{
		Jar: jar,
	}}
	c.client.CheckRedirect = c.saveRedirect
	return c, nil
}

func (c *cli) get(url string) (io.ReadCloser, error) {
	resp, err := c.client.Get(url) // не очень понятно, в случае 503 вернется ошибка или статус код надо проверять?
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, rejected
	}
	return resp.Body, err
}

func (c *cli) post(data url.Values) (url.Values, error) {
	resp, err := c.client.PostForm(c.url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, rejected
	}
	return parseHtml(resp.Body)
}

func (c *cli) parseInitialPage() error {
	r, err := c.get(mainUrl)
	if err != nil {
		return err
	}
	defer r.Close()
	link, err := findQuestion1Link(r)
	if err != nil {
		return err
	}
	if link == "" {
		return errors.New("not found")
	}
	c.url = mainUrl + link
	return nil
}

func (c *cli) parseNextPages() error {
	r, err := c.get(c.url)
	if err != nil {
		return err
	}
	defer r.Close()
	for data, err := parseHtml(r); err == nil && data != nil; data, err = c.post(data) {
	}
	return err
}

func main() {
	cli, err := newCli()
	if err != nil {
		log.Fatal(err)
	}
	if err := cli.parseInitialPage(); err != nil {
		log.Fatal(err)
	}
	if err := cli.parseNextPages(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}
