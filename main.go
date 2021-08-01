package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
)

var rejected = errors.New("rejected")

type cli struct {
	client *http.Client // вопрос, а может *http.Client
}

func newCli() (*cli, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	return &cli{client: &http.Client{
		Jar: jar,
	}}, nil
}

func (c *cli) get(url string) (io.ReadCloser, error) {
	resp, err := c.client.Get(url) // не очень понятно, в случае 503 вернется ошибка или статус код надо проверять?
	if err != nil {
		return nil, err
	} else if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, rejected
	}
	return resp.Body, nil
}

func (c *cli) req1() error {
	body, err := c.get("http://test.youplace.net/")
	if err != nil {
		return err
	}
	defer body.Close()
	return err
}

func (c *cli) req2() error {
	body, err := c.get("http://test.youplace.net/question/1")
	if err != nil {
		return err
	}
	defer body.Close()

	var b strings.Builder
	if _, err = io.Copy(&b, body); err != nil {
		return err
	}
	s := b.String()

	fmt.Println(s)

	var data url.Values

	resp, err := c.client.PostForm("http://test.youplace.net/question/1", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		return err
	}
	return nil
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
}
