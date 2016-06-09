package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/mvdan/xurls"
	"github.com/thoj/go-ircevent"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func main() {
	c := irc.IRC("luggy", "luggy")
	c.AddCallback("001", func(e *irc.Event) {
		c.Join("#utdlug")
	})
	c.AddCallback("PRIVMSG", func(e *irc.Event) {
		msg := e.Message()
		if e.Nick == "lugbrahtomy" && strings.Contains(msg, "wikipedia") {
			return
		}
		for _, link := range xurls.Relaxed.FindAllString(msg, -1) {
			u, err := url.Parse(link)
			if err != nil {
				continue
			}
			if !strings.HasPrefix(u.Scheme, "http") {
				link = "http://" + link
			}
			res, err := http.Get(link)
			if err != nil {
				continue
			}
			root, err := html.Parse(res.Body)
			defer res.Body.Close()
			if err != nil {
				continue
			}
			title, ok := scrape.Find(root, scrape.ByTag(atom.Title))
			if ok {
				c.Privmsg(e.Arguments[0], scrape.Text(title))
			}
		}
	})
	if err := c.Connect("irc.oftc.net:6667"); err != nil {
		panic(err)
	}
	c.Loop()
}
