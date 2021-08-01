package main

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestForm(t *testing.T) {
	resp := `
	<form method="POST">

    <p><input type="radio" name="UBuBQIWEamvB4aXp" value="tqgBihtyp5">tqgBihtyp5 &nbsp;&nbsp;&nbsp;<input type="radio" name="UBuBQIWEamvB4aXp" value="U5z7xFhMsr">U5z7xFhMsr &nbsp;&nbsp;&nbsp;<input type="radio" name="UBuBQIWEamvB4aXp" value="9i4lVrAK5">9i4lVrAK5 &nbsp;&nbsp;&nbsp;<input type="radio" name="UBuBQIWEamvB4aXp" value="pMeJe2q">pMeJe2q &nbsp;&nbsp;&nbsp;<input type="radio" name="UBuBQIWEamvB4aXp" value="91sAV9LJEgYapcvc">91sAV9LJEgYapcvc &nbsp;&nbsp;&nbsp;</p>

    <p><select name="BOb1kMQ8T9lu7NPt"><option value=""></option><option value="i8jmbynhPf">i8jmbynhPf</option><option value="5S7icQmd">5S7icQmd</option><option value="irHSUXCu">irHSUXCu</option><option value="y4QgycIuqoFvxyF">y4QgycIuqoFvxyF</option><option value="NloFRczkt">NloFRczkt</option></select></p>

    <p><input type="text" name="OZB2YxAEdU7hXiaZ"></p>

    <p><input type="text" name="ajivmQlpOk1s8O3T"></p>

    <p><button type="submit">Submit</button></p>
	</form>
	`
	data, _ := parseHtml(strings.NewReader(resp))
	if data.Get("ajivmQlpOk1s8O3T") != "test" {
		t.Error("ajivmQlpOk1s8O3T")
	}
	if data.Get("OZB2YxAEdU7hXiaZ") != "test" {
		t.Error("OZB2YxAEdU7hXiaZ")
	}
	if data.Get("UBuBQIWEamvB4aXp") != "91sAV9LJEgYapcvc" {
		t.Error("UBuBQIWEamvB4aXp")
	}
	if data.Get("BOb1kMQ8T9lu7NPt") != "y4QgycIuqoFvxyF" {
		t.Error("BOb1kMQ8T9lu7NPt")
	}
	t.Log(data)
}

func TestHtml(t *testing.T) {
	r := strings.NewReader(`<html><a href="link">test</a></html>`)
	z := html.NewTokenizer(r)
	for {
		tt := z.Next()
		if z.Err() != nil {
			t.Log(z.Err())
			return
		}
		tn := z.Token()
		if tt == html.StartTagToken && tn.Data == "a" {
			for _, a := range tn.Attr {
				if a.Key == "href" {
					t.Log(a.Namespace, a.Val)
				}
			}

		}

	}
}
