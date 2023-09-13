package main

import (
	"fmt" // пакет для форматированного ввода вывода
	"github.com/mmcdole/gofeed"
	"log" // пакет для логирования
	"net/http"
)

func getMainHtml() string {
	return "<!DOCTYPE html>\n<html lang=\"en\">\n\n<head>\n  " +
		"<meta charset=\"UTF-8\" />\n  <meta name=\"viewport\" " +
		"content=\"width=device-width, initial-scale=1.0\" />\n  " +
		"<link rel=\"stylesheet\" href=\"style.css\" />\n  " +
		"<title>Browser</title>\n</head>\n\n<body>\n  <h1>\n    " +
		"Main menu\n  </h1>\n  <ul>\n  <li><a href=\"/rss\">RSS</a></li>\n  <li>Tea</li>\n  " +
		"<li>Milk</li>\n</ul>\n</body>\n\n</html>"
}

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, getMainHtml()) // отправляем данные на клиентскую сторону
}

func RssRouterHandler(w http.ResponseWriter, r *http.Request) {
	getRss, _ := http.Get("https://www.onliner.by/feed")
	defer getRss.Body.Close()
	feedParser := gofeed.NewParser()
	feed, _ := feedParser.Parse(getRss.Body)
	//fmt.Println(feed.Title)
	responseHtml := "<!DOCTYPE html>\n" +
		"<html lang=\"en\">\n\n<head>\n  <meta charset=\"UTF-8\" />\n" +
		"  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\" />\n" +
		"  <link rel=\"stylesheet\" href=\"style.css\" />\n" +
		"  <title>Browser</title>\n</head>\n\n<body>\n"
	responseHtml += "<h1><a href=\"" + feed.Link + "\">" + feed.Title + "</a></h1>"
	for _, v := range feed.Items {
		responseHtml += v.Description
	}
	responseHtml += "</body>"
	fmt.Fprintf(w, "%s", responseHtml) // отправляем данные на клиентскую сторону
}

func main() {
	http.HandleFunc("/", HomeRouterHandler) // установим роутер
	http.HandleFunc("/rss", RssRouterHandler)
	err := http.ListenAndServe(":9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
