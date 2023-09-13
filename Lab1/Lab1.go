package main

import (
	"bufio"
	"fmt" // пакет для форматированного ввода вывода
	"github.com/mmcdole/gofeed"
	"html/template"
	"log" // пакет для логирования
	"net/http"
	"os"
)

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("Lab1/main.html")
	tmpl.Execute(w, nil)
}

func RssRouterHandler(w http.ResponseWriter, r *http.Request) {
	//interfaxURL := "https://www.interfax.ru/rss.asp"
	onlinerURL := "https://www.onliner.by/feed"
	getRss, _ := http.Get(onlinerURL)
	defer getRss.Body.Close()
	feedParser := gofeed.NewParser()
	feed, _ := feedParser.Parse(getRss.Body)

	response := ""
	f, _ := os.Open("Lab1/rss.html")
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		response += scanner.Text()
	}
	//fmt.Println(response)
	response += "<h1><a href=\"" + feed.Link + "\">" + feed.Title + "</a></h1>"
	for _, v := range feed.Items {
		//response += "<p><a href =\"" + v.Link + "\">" + v.Description + "</a></p>"
		response += v.Description
	}
	response += "</body></html>"
	fmt.Fprintf(w, "%s", response) // отправляем данные на клиентскую сторону
}

func PostRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if len(r.Form) != 0 {
		fmt.Println(r.Form)
		fmt.Fprintf(w, "You are from "+r.Form.Get("to"))
	} else {
		tmpl, _ := template.ParseFiles("Lab1/post.html")
		tmpl.Execute(w, nil) // отправляем данные на клиентскую сторону
	}

}

func main() {
	http.HandleFunc("/", HomeRouterHandler) // установим роутер
	http.HandleFunc("/rss", RssRouterHandler)
	http.HandleFunc("/post", PostRouterHandler)
	err := http.ListenAndServe(":9000", nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
