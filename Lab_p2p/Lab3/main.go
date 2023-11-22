package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"

	p2p "github.com/leprosus/golang-p2p"
)

type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type NLogger struct{}

func NewStdLogger() (l *NLogger) {
	return &NLogger{}
}

func (l *NLogger) Info(msg string) {

}

func (l *NLogger) Warn(msg string) {

}

func (l *NLogger) Error(msg string) {

}

type Book struct {
	Values []string
}

type Node struct {
	Port       int
	Parent     int
	NumOfChild int
	Children   []int
	Book       Book
	Version    int
}

func sendHelloParent(this *Node) {
	tcp := p2p.NewTCP("localhost", strconv.Itoa(this.Parent))
	client, err := p2p.NewClient(tcp)
	client.SetLogger(NewStdLogger())
	if err != nil {
		log.Panicln(err)
	}

	req := p2p.Data{}
	err = req.SetJson(this)
	if err != nil {
		log.Panicln(err)
	}

	//_ := p2p.Data{}
	got := Node{}

	res, _ := client.Send("AddNewChild", req)
	res.GetJson(&got)
	this.Book.Values = got.Book.Values

}

func startServer(this *Node) {
	tcp := p2p.NewTCP("localhost", strconv.Itoa(this.Port))

	server, err := p2p.NewServer(tcp)
	if err != nil {
		log.Panicln(err)
	}
	server.SetLogger(NewStdLogger())

	server.SetHandle("AddNewChild", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		childInfo := Node{}
		err = req.GetJson(&childInfo)
		if err != nil {
			return
		}
		this.Children = append(this.Children, childInfo.Port)
		//fmt.Println(this.Children)

		res = p2p.Data{}
		err = res.SetJson(this)

		return
	})

	server.SetHandle("UpdateBook", func(ctx context.Context, req p2p.Data) (res p2p.Data, err error) {
		Info := Node{}
		err = req.GetJson(&Info)
		if reflect.DeepEqual(Info.Book.Values, this.Book.Values) {
			return
		}
		if err != nil {
			return
		}
		this.Book = Info.Book
		this.Version = Info.Version
		if this.Parent != 0 {
			tcp := p2p.NewTCP("localhost", strconv.Itoa(this.Parent))
			client, err := p2p.NewClient(tcp)
			client.SetLogger(NewStdLogger())
			if err != nil {
				log.Panicln(err)
			}

			req := p2p.Data{}
			err = req.SetJson(this)
			if err != nil {
				log.Panicln(err)
			}
			client.Send("UpdateBook", req)
		}

		for i := 0; i < len(this.Children); i++ {
			tcp := p2p.NewTCP("localhost", strconv.Itoa(this.Children[i]))
			client, err := p2p.NewClient(tcp)
			client.SetLogger(NewStdLogger())
			if err != nil {
				log.Panicln(err)
			}

			req := p2p.Data{}
			err = req.SetJson(this)
			if err != nil {
				log.Panicln(err)
			}
			client.Send("UpdateBook", req)
		}

		//fmt.Println(this.Book.Values)
		return
	})

	err = server.Serve()
	if err != nil {
		log.Panicln(err)
	}
}

func startbook(this *Node) {
	for {
		command := ""
		fmt.Print("enter command: ")
		fmt.Scan(&command)
		if command == "add" {
			val := ""
			fmt.Print("enter value: ")
			fmt.Scan(&val)
			this.Book.Values = append(this.Book.Values, val)
			this.Version++
			tcp := p2p.NewTCP("localhost", strconv.Itoa(this.Parent))
			client, err := p2p.NewClient(tcp)
			client.SetLogger(NewStdLogger())
			if err != nil {
				log.Panicln(err)
			}

			req := p2p.Data{}
			err = req.SetJson(this)
			if err != nil {
				log.Panicln(err)
			}
			client.Send("UpdateBook", req)
		} else if command == "get" {
			fmt.Println(this.Book)
		} else if command == "delete" {
			val := ""
			fmt.Print("what to delete?: ")
			fmt.Scan(&val)
			for i := 0; i < len(this.Book.Values); i++ {
				if this.Book.Values[i] == val {
					this.Book.Values = append(this.Book.Values[:i], this.Book.Values[i+1:]...)
				}
			}
			this.Version++
			tcp := p2p.NewTCP("localhost", strconv.Itoa(this.Parent))
			client, err := p2p.NewClient(tcp)
			client.SetLogger(NewStdLogger())
			if err != nil {
				log.Panicln(err)
			}

			req := p2p.Data{}
			err = req.SetJson(this)
			if err != nil {
				log.Panicln(err)
			}
			client.Send("UpdateBook", req)
		}
	}

}

func main() {
	this := Node{
		Port:       0,
		Parent:     0,
		NumOfChild: 0,
		Children:   make([]int, 0),
		Book:       Book{Values: make([]string, 0)},
		Version:    0,
	}

	fmt.Println("Port: ")
	fmt.Scan(&this.Port)
	fmt.Println("Parent: ")

	fmt.Scan(&this.Parent)

	if this.Parent != 0 {
		go sendHelloParent(&this)
	}
	go startbook(&this)
	startServer(&this)

}
