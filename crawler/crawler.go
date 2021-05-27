package crawler

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gocolly/colly/v2"
)

type Walk struct {
	Uri  string
	Next string
	Item string

	Cookies []*http.Cookie

	LimitItems int
	itemsCount int
	LimitPages int
	pagesCount int
}

type Element struct {
	*colly.HTMLElement
}

func (w *Walk) Start(fn func(e *Element)) error {
	if w.Uri == "" || w.Next == "" || w.Item == "" {
		return errors.New("uri/href/item not defined")
	}

	c := colly.NewCollector()

	w.setCookies(c)

	c.OnHTML(w.Item, func(e *colly.HTMLElement) {
		if w.reachLimit() {
			return
		}
		fn(&Element{e})
		w.itemsCount += 1
	})

	c.OnHTML(w.Next, func(e *colly.HTMLElement) {
		if w.reachLimit() {
			return
		}
		w.pagesCount += 1
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.Visit(w.Uri)

	return nil
}

func (w *Walk) reachLimit() bool {
	return (w.LimitItems != 0 && w.itemsCount >= w.LimitItems) || (w.LimitPages != 0 && w.pagesCount >= w.LimitPages)
}

func (w *Walk) setCookies(c *colly.Collector) {
	if w.Cookies == nil {
		return
	}
	u, _ := url.Parse(w.Uri)
	u.Path = ""
	host := u.String()
	c.SetCookies(host, w.Cookies)
}
