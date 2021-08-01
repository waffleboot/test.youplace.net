package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

var rejected = errors.New("rejected")

type cli struct {
	client *http.Client // вопрос, а может *http.Client
	url    string
}

func (c *cli) checkRedirect(req *http.Request, via []*http.Request) error {
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
	c.client.CheckRedirect = c.checkRedirect
	return c, nil
}

func (c *cli) get(url string) (string, error) {
	resp, err := c.client.Get(url) // не очень понятно, в случае 503 вернется ошибка или статус код надо проверять?
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusServiceUnavailable {
		return "", rejected
	}
	var b strings.Builder
	_, err = io.Copy(&b, resp.Body)
	return b.String(), err
}

func (c *cli) req1() (err error) {
	resp, err := c.get("http://test.youplace.net/")
	link, err := form3(resp)
	if err != nil {
		return err
	}
	c.url = "http://test.youplace.net" + link
	return
}

func (c *cli) req2() error {
	resp, err := c.get(c.url)
	if err != nil {
		return err
	}
	data, err := form1(resp)
	if err != nil {
		if err != io.EOF {
			return err
		}
		return nil
	}
	for {
		data, err = c.req3(data)
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		<-time.After(1 * time.Second)
	}
}

func (c *cli) req3(data url.Values) (url.Values, error) {
	resp, err := c.client.PostForm(c.url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return form2(resp.Body)
}

func (c *cli) run() (err error) {
	for _, m := range []func() error{c.req1, c.req2} {
		if err = m(); err != nil {
			return
		}
	}
	return
}

func main() {
	c, err := newCli()
	if err != nil {
		log.Fatal(err)
	}
	if err := c.run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("done")
}
