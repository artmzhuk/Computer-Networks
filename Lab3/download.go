package main

import (
	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Ref, Time, Title string
	val              int64
}

func readItem(item *html.Node) *Item {
	if a := item.FirstChild; isElem(a, "a") {
		if cs := getChildren(a); len(cs) == 2 && isElem(cs[0], "time") && isText(cs[1]) {
			return &Item{
				Ref:   getAttr(a, "href"),
				Time:  getAttr(cs[0], "title"),
				Title: cs[1].Data,
			}
		}
	}
	return nil
}

func search(node *html.Node) []*Item {
	if isDiv(node, "sc-66133f36-2 cgmess") {
		counter := 0
		var items []*Item
		for c := node.FirstChild.FirstChild.NextSibling.NextSibling; c != nil; c = c.NextSibling {
			if isElem(c, "tbody") {
				for ch := c.FirstChild; ch != nil; ch = ch.NextSibling {
					if isElem(ch, "tr") {
						currency := ch.FirstChild.NextSibling.NextSibling.FirstChild.FirstChild

						if currency.Data != "span" {
							currCap := ch.LastChild.PrevSibling.PrevSibling.FirstChild.FirstChild.FirstChild.FirstChild.Data
							url := getAttr(currency, "href")
							currName := currency.FirstChild.LastChild.FirstChild.FirstChild.FirstChild.Data
							//t := currency.FirstChild.Data
							val, _ := strconv.ParseInt(strings.Replace(currCap[1:], ",", "", -1), 10, 64)
							items = append(items, &Item{
								Ref:   url,
								Time:  currCap,
								Title: currName,
								val:   val,
							})
						} else if counter < 99 {
							url := getAttr(currency.Parent, "href")
							currName := currency.Parent.FirstChild.NextSibling.FirstChild.Data
							//t := currency.FirstChild.Data
							items = append(items, &Item{
								Ref:   url,
								Time:  "¯\\_(ツ)_/¯",
								Title: currName,
								val:   0,
							})
						}
						counter++
					}
				}
			}
		}
		sort.SliceStable(items, func(i, j int) bool {
			return items[i].val > items[j].val
		})
		return items
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

func downloadNews() []*Item {
	log.Info("sending request ")
	if response, err := http.Get("https://coinmarketcap.com/"); err != nil {
		log.Error("request failed", "error", err)
	} else {
		//defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response ", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from ", "error", err)
			} else {
				log.Info("HTML  parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
