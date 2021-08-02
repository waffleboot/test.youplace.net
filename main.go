package main

import (
	"errors"
	"fmt"
	"log"
)

const mainUrl = "http://test.youplace.net"

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
	c.nextPageUrl = mainUrl + link
	return nil
}

func (c *cli) parseNextPages() error {
	r, err := c.get(c.nextPageUrl)
	if err != nil {
		return err
	}
	defer r.Close()
	for data, err := parseHtml(r); err == nil && data != nil; data, err = c.post(data) {
	}
	return err
}
