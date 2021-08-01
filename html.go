package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

var name = regexp.MustCompile(`name="([^"]+)"`)

var input = regexp.MustCompile(`<input[^>]+>`)

var value = regexp.MustCompile(`value="([^"]+)"`)

var gselect = regexp.MustCompile(`<select[^>]+>.*</select>`)

var option = regexp.MustCompile(`<option[^>]+>`)

var passed = regexp.MustCompile(`Test successfully passed`)

func findQuestion1Link(resp string) (string, error) {
	var ans string
	r := strings.NewReader(resp)
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if err := z.Err(); err != nil {
			if err != io.EOF {
				return "", err
			}
			return ans, nil
		}
		if tn := z.Token(); tt == html.StartTagToken && tn.Data == "a" {
			for _, a := range tn.Attr {
				if a.Key == "href" {
					ans = a.Val
				}
			}

		}
	}
}

func parseString(resp string) url.Values {
	if passed.FindString(resp) != "" {
		return nil
	}
	data := url.Values{}
	radio := make(map[string]string)
	for _, s := range input.FindAllString(resp, -1) {
		n := name.FindStringSubmatch(s)[1]
		if strings.Contains(s, `type="text"`) {
			data.Add(n, "test")
		} else if strings.Contains(s, `type="radio"`) {
			p := radio[n]
			v := value.FindStringSubmatch(s)[1]
			if len(v) > len(p) {
				radio[n] = v
			}
		}
	}
	for _, s := range gselect.FindAllString(resp, -1) {
		n := name.FindStringSubmatch(s)[1]
		for _, o := range option.FindAllString(s, -1)[1:] {
			p := radio[n]
			v := value.FindStringSubmatch(o)[1]
			if len(v) > len(p) {
				radio[n] = v
			}
		}
	}
	for k, v := range radio {
		data.Add(k, v)
	}
	return data
}

func parseReader(resp io.Reader) (data url.Values, err error) {
	var b strings.Builder
	if _, err = io.Copy(&b, resp); err != nil {
		return
	}
	fmt.Fprintln(os.Stderr, b.String())
	data = parseString(b.String())
	return
}
