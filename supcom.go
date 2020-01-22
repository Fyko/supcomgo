package supcmgo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Votes item votes
type Votes struct {
	Positive int
	Negative int
}

// Prices item prices
type Prices struct {
	Combined string
	USD      int
	GBP      int
}

// Item a Supreme Community item instance
type Item struct {
	Name        string
	Price       Prices
	Image       string
	Description string
	Category    string
	Votes       Votes
}

// Items all the items
type Items []Item

// FetchLatest fetch the latest supreme community droplist
func FetchLatest() string {
	c := colly.NewCollector(
		colly.AllowedDomains("www.supremecommunity.com"),
		colly.Async(true),
	)

	var link string

	c.OnHTML(".block", func(e *colly.HTMLElement) {
		link = e.Attr("href")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("OnError ", r.StatusCode, err)
	})

	c.Visit("https://www.supremecommunity.com/season/latest/droplists/")
	c.Wait()

	return link
}

// FetchDroplist fetches and returns the droplist from the provided URL
func FetchDroplist(url string) []Item {
	c := colly.NewCollector(
		colly.AllowedDomains("www.supremecommunity.com"),
		colly.Async(true),
	)

	var items Items

	c.OnHTML(".masonry__item", func(e *colly.HTMLElement) {
		var item Item
		item.Name = e.ChildAttr(".card-details", "data-itemname")

		price := e.ChildText(".label-price")
		priceparts := strings.Split(price, "/")
		item.Price.Combined = price
		item.Price.USD, _ = strconv.Atoi(trimLeftChar(priceparts[0]))
		item.Price.GBP, _ = strconv.Atoi(trimLeftChar(priceparts[len(priceparts)-1]))

		item.Image = "https://www.supremecommunity.com" + e.ChildAttr(".prefill-img", "src")
		description := e.ChildAttr("img", "alt")

		parts := strings.Split(description, "- ")
		item.Description = parts[len(parts)-1]

		item.Votes.Negative, _ = strconv.Atoi(e.ChildText(".progress-bar-danger.droplist-vote-bar"))
		item.Votes.Positive, _ = strconv.Atoi(e.ChildText(".progress-bar-success.droplist-vote-bar"))

		item.Category = e.ChildText(".category")

		items = append(items, item)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("OnError ", r.StatusCode, err)
	})

	c.Visit("https://www.supremecommunity.com/" + url)
	c.Wait()

	return items
}

func trimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}
