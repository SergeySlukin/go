package main

import (
	"github.com/opesun/goquery"
	"strings"
)

func Grab(g chan<- string) {
	x, err := goquery.ParseUrl("http://vpustotu.ru/moderation/")
	if err == nil {
		if s := strings.TrimSpace(x.Find(".fi_text").Text()); s != "" {
			g <- s
		}
	} else {
		g <- err.Error()
	}


}
