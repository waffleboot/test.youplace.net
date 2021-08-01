package main

import (
	"errors"
	"io"
	"net/url"
	"regexp"
	"strings"
)

var name = regexp.MustCompile(`name="([^"]+)"`)

var input = regexp.MustCompile(`<input[^>]+>`)

var value = regexp.MustCompile(`value="([^"]+)"`)

var gselect = regexp.MustCompile(`<select[^>]+>.*</select>`)

var option = regexp.MustCompile(`<option[^>]+>`)

var passed = regexp.MustCompile(`Test successfully passed`)

var start = regexp.MustCompile(`<a href="(/question/1)"><button>Start test</button></a>`)

func form3(resp string) (string, error) {
	for _, link := range start.FindStringSubmatch(resp)[1:] {
		return link, nil
	}
	return "", errors.New("not found")
}

func form1(resp string) (url.Values, error) {
	if passed.FindString(resp) != "" {
		return nil, io.EOF
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
	return data, nil
}

func form2(resp io.Reader) (url.Values, error) {
	var b strings.Builder
	if _, err := io.Copy(&b, resp); err != nil {
		return url.Values{}, err
	}
	return form1(b.String())
}
