package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

func findQuestion1Link(r io.Reader) (ans string, err error) {
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if err = z.Err(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
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

func ntv(tn html.Token) (n string, t string, v string) {
	for _, a := range tn.Attr {
		switch a.Key {
		case "name":
			n = a.Val
		case "type":
			t = a.Val
		case "value":
			v = a.Val
		}
	}
	return
}

func parseHtml(r io.Reader) (data url.Values, err error) {
	data = url.Values{}
	z := html.NewTokenizer(r)
	var gn string
	for {
		tt := z.Next()
		if err = z.Err(); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		tn := z.Token()
		if tt == html.StartTagToken {
			switch tn.Data {
			case "input":
				fmt.Fprintln(os.Stderr, tn)
				n, t, v := ntv(tn)
				if t == "text" {
					data.Add(n, "test")
				} else if t == "radio" {
					if p := data.Get(n); len(v) > len(p) {
						data.Set(n, v)
					}
				}
			case "select":
				fmt.Fprintln(os.Stderr, tn)
				gn, _, _ = ntv(tn)
			case "option":
				fmt.Fprintln(os.Stderr, tn)
				_, _, v := ntv(tn)
				if p := data.Get(gn); len(v) > len(p) {
					data.Set(gn, v)
				}
			}
		} else if tt == html.TextToken && tn.Data == "Test successfully passed" {
			return nil, nil
		}
	}
}
